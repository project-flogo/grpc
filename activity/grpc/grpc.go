package grpc

import (
	"errors"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func init() {
	activity.Register(&Activity{}, New)
}

type mashGRPCClienConn struct {
	conn  *grpc.ClientConn
	count int
}

type mashGRPCClinetConns struct {
	connMap map[string]*mashGRPCClienConn
	sync.Mutex
}

var conns = mashGRPCClinetConns{
	connMap: make(map[string]*mashGRPCClienConn),
}

// Activity is a GRPC activity
type Activity struct {
	settings *Settings
}

// New creates a new javascript activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := Settings{}
	err := metadata.MapToStruct(ctx.Settings(), &settings, true)
	if err != nil {
		return nil, err
	}

	logger := ctx.Logger()
	logger.Debugf("Setting: %b", settings)

	act := Activity{
		settings: &settings,
	}

	return &act, nil
}

// Metadata return the metadata for the activity
func (a *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the activity
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()

	input := Input{}
	err = ctx.GetInputObject(&input)
	if err != nil {
		return false, err
	}

	opts := []grpc.DialOption{}
	logger.Debug("enableTLS: ", a.settings.EnableTLS)
	if a.settings.EnableTLS {
		logger.Debug("ClientCert: ", a.settings.ClientCert)
		creds, err := credentials.NewClientTLSFromFile(a.settings.ClientCert, "")
		if err != nil {
			logger.Error(err)
		}

		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	} else {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	}

	conn, err := getConnection(a.settings.HostURL, logger, opts)
	if err != nil {
		return false, err
	}
	defer releaseConnection(a.settings.HostURL)

	logger.Debug("operating mode: ", a.settings.OperatingMode)

	output := Output{}
	switch a.settings.OperatingMode {
	case "grpc-to-grpc":
		err = a.gRPCTogRPCHandler(&input, &output, logger, conn)
		if err != nil {
			return false, err
		}
		err = ctx.SetOutputObject(&output)
		if err != nil {
			return false, err
		}
		return true, nil
	case "rest-to-grpc":
		err = a.restTogRPCHandler(&input, &output, logger, conn)
		if err != nil {
			return false, err
		}
		err = ctx.SetOutputObject(&output)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	logger.Error(errors.New("Invalid use of service , OperatingMode not recognised"))
	return false, errors.New("Invalid use of service , OperatingMode not recognised")
}

// getconnection returns single client connection object per hostaddress
func getConnection(hostAdds string, logger log.Logger, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	conns.Lock()
	defer conns.Unlock()
	conn := conns.connMap[hostAdds]
	if conn == nil {
		c, err := grpc.Dial(hostAdds, opts...)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		conn = &mashGRPCClienConn{
			conn:  c,
			count: 0,
		}
		conns.connMap[hostAdds] = conn
	}
	conn.count++
	return conn.conn, nil
}

// releaseConnection closes created client connection per hostaddress
func releaseConnection(hostAdds string) {
	conns.Lock()
	defer conns.Unlock()
	conn := conns.connMap[hostAdds]
	conn.count--
	if conn.count <= 0 {
		conn.conn.Close()
		delete(conns.connMap, hostAdds)
	}
}
