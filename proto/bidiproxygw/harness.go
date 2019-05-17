package bidiproxygw

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rameshpolishetti/movingavg/cma"
	"github.com/rameshpolishetti/movingavg/sma"
	"google.golang.org/grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
	Server       *grpc.Server
	payladMapAry map[string]Payload
}

var clntRcvdCount []int64

var movingRspTm *sma.ThreadSafeSMA
var cumulativeRspTm *cma.ThreadSafeCMA
var rspTmAvg []float64
var tPSAvg []int64

var clntEOFCount []int

var exitSignal bool

var totalSc, totalRc int64

var minRT, maxRT int64 = 1000000000, 0

func CallClient(port *string, threads int) (interface{}, error) {
	hostIP := os.Getenv("HOSTIP")
	if len(hostIP) == 0 {
		hostIP = "localhost"
	}

	clientAddr := *port
	if len(*port) == 0 {
		clientAddr = hostIP + ":9000"
	} else {
		clientAddr = hostIP + ":" + *port
	}

	conn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewSampleServiceBidiClient(conn)

	ma, _ := sma.New(600000)
	movingRspTm = sma.GetThreadSafeSMA(ma)

	ma1 := cma.New()
	cumulativeRspTm = cma.GetThreadSafeCMA(ma1)

	clntRcvdCount = make([]int64, threads)

	clntEOFCount = make([]int, threads)

	exitSignal = false

	fmt.Println("Starting threads", time.Now())

	go responseTime(threads)

	for thread := 0; thread < threads; thread++ {
		time.Sleep(20 * time.Millisecond)
		fmt.Println(thread)
		go bidiStreaming(client, thread)
	}
	fmt.Println("All Threads Activated", time.Now())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	fmt.Println("Exit signal received", time.Now())
	exitSignal = true

	for {
		time.Sleep(time.Second)
		count := 1
		for i := 0; i < threads; i++ {
			if clntEOFCount[i] == 1 {
				count++
			} else {
				count = 1
				break
			}
			fmt.Println("count: ", count, "threads", threads)
			if count == threads {
				fmt.Println("All threads completed")
				break
			}
		}
		if count != 1 {
			fmt.Println("All threads completed")

			fmt.Println("Threads:", threads)
			fmt.Println("TimeBtwnMessages:30Ms")
			fmt.Println("TimeBtwnThread:20Ms")
			fmt.Println("Atomic send count: ", totalSc)
			fmt.Println("Atomic receive count: ", totalRc)
			fmt.Println("Response Samples: ", rspTmAvg)
			fmt.Println("TPS Samples: ", tPSAvg)
			fmt.Println("AverageResponseTime:", cumulativeRspTm.Avg())
			fmt.Println("slow response time: ", maxRT)
			fmt.Println("fast response time: ", minRT)

			return nil, nil
		}
	}

}

func responseTime(threads int) {
	tick := time.Tick(6 * time.Second)
	var totalCrc int64

	for {
		select {
		case <-tick:
			avg := movingRspTm.Avg()
			rspTmAvg = append(rspTmAvg, avg/float64(1000000))

			tPSAvg = append(tPSAvg, totalRc-totalCrc)

			//fmt.Println("Tps=", totalRc-totalCrc)

			totalCrc = totalRc
		}
	}
}

func totalCount(arr []int) int {
	total := 0
	for i := 0; i < len(arr); i++ {
		total = total + arr[i]
	}
	return total
}

func bidiStreaming(client SampleServiceBidiClient, thread int) error {
	stream, err := client.BidiStreaming(context.Background())
	if err != nil {
		fmt.Println("err1", err)
		return err
	}
	var rspTm, ttlRspTm int64

	waitc := make(chan struct{})
	go func() {
		for {
			payload, err := stream.Recv()
			if err == io.EOF {
				clntRcvdCount[thread] = ttlRspTm
				clntEOFCount[thread] = 1
				close(waitc)

				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a payload : %v", err)
				close(waitc)
				return
			}
			if payload != nil {
				atomic.AddInt64(&totalRc, 1)
				t1 := time.Now().UnixNano()
				rspTm = t1 - payload.GetTimestamp1()
				movingRspTm.AddSampleInt64(rspTm)
				cumulativeRspTm.AddSampleInt64(rspTm)
				if rspTm > 1000000000 {
					fmt.Println("Response Time:", rspTm, time.Now())
				}
				if rspTm < minRT {
					atomic.StoreInt64(&minRT, rspTm)
				}

				if rspTm > maxRT {
					atomic.StoreInt64(&maxRT, rspTm)
				}
			}
		}
	}()
	for {
		tN := time.Now().UnixNano()
		payload := Payload{
			Timestamp1: tN,
		}
		if err := stream.Send(&payload); err != nil {
			fmt.Println("error while sending payload", payload, err)
			return err
		}
		atomic.AddInt64(&totalSc, 1)
		time.Sleep(10 * time.Millisecond)
		if exitSignal {
			fmt.Println("Send Close signal recieved")
			break
		}
	}
	fmt.Println("Closing for thread: ", thread)
	stream.CloseSend()
	<-waitc

	return nil
}

// CallServer registers and calls sample bidi server
func CallServer() (*ServerStrct, error) {
	s := grpc.NewServer()
	server := ServerStrct{
		Server:       s,
		payladMapAry: make(map[string]Payload),
	}

	RegisterSampleServiceBidiServer(s, &server)

	return &server, nil
}

// Start creates running grpc instance
func (t *ServerStrct) Start() error {
	addr := ":9000"
	log.Println("Starting server on port: ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error occures while starting server: ", err)
		return err
	}

	return t.Server.Serve(lis)
}

func (t *ServerStrct) BidiStreaming(bReq SampleServiceBidi_BidiStreamingServer) error {
	fmt.Println("server BidiStreaming method called")
	for {
		payload, err := bReq.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error occured while receiving", err)
			break
		}

		if err := bReq.Send(payload); err != nil {
			fmt.Println("error while sending payload", payload)
			return err
		}
	}
	return nil
}
