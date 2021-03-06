// Code generated by protoc-gen-go. DO NOT EDIT.
// source: petstore.proto

package rest2grpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Pet struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pet) Reset()         { *m = Pet{} }
func (m *Pet) String() string { return proto.CompactTextString(m) }
func (*Pet) ProtoMessage()    {}
func (*Pet) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{0}
}
func (m *Pet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pet.Unmarshal(m, b)
}
func (m *Pet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pet.Marshal(b, m, deterministic)
}
func (dst *Pet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pet.Merge(dst, src)
}
func (m *Pet) XXX_Size() int {
	return xxx_messageInfo_Pet.Size(m)
}
func (m *Pet) XXX_DiscardUnknown() {
	xxx_messageInfo_Pet.DiscardUnknown(m)
}

var xxx_messageInfo_Pet proto.InternalMessageInfo

func (m *Pet) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Pet) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type User struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email                string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{1}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type PetByIdRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PetByIdRequest) Reset()         { *m = PetByIdRequest{} }
func (m *PetByIdRequest) String() string { return proto.CompactTextString(m) }
func (*PetByIdRequest) ProtoMessage()    {}
func (*PetByIdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{2}
}
func (m *PetByIdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PetByIdRequest.Unmarshal(m, b)
}
func (m *PetByIdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PetByIdRequest.Marshal(b, m, deterministic)
}
func (dst *PetByIdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PetByIdRequest.Merge(dst, src)
}
func (m *PetByIdRequest) XXX_Size() int {
	return xxx_messageInfo_PetByIdRequest.Size(m)
}
func (m *PetByIdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PetByIdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PetByIdRequest proto.InternalMessageInfo

func (m *PetByIdRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type UserByNameRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserByNameRequest) Reset()         { *m = UserByNameRequest{} }
func (m *UserByNameRequest) String() string { return proto.CompactTextString(m) }
func (*UserByNameRequest) ProtoMessage()    {}
func (*UserByNameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{3}
}
func (m *UserByNameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserByNameRequest.Unmarshal(m, b)
}
func (m *UserByNameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserByNameRequest.Marshal(b, m, deterministic)
}
func (dst *UserByNameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserByNameRequest.Merge(dst, src)
}
func (m *UserByNameRequest) XXX_Size() int {
	return xxx_messageInfo_UserByNameRequest.Size(m)
}
func (m *UserByNameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserByNameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserByNameRequest proto.InternalMessageInfo

func (m *UserByNameRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

type PetResponse struct {
	Pet                  *Pet     `protobuf:"bytes,1,opt,name=pet,proto3" json:"pet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PetResponse) Reset()         { *m = PetResponse{} }
func (m *PetResponse) String() string { return proto.CompactTextString(m) }
func (*PetResponse) ProtoMessage()    {}
func (*PetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{4}
}
func (m *PetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PetResponse.Unmarshal(m, b)
}
func (m *PetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PetResponse.Marshal(b, m, deterministic)
}
func (dst *PetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PetResponse.Merge(dst, src)
}
func (m *PetResponse) XXX_Size() int {
	return xxx_messageInfo_PetResponse.Size(m)
}
func (m *PetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PetResponse proto.InternalMessageInfo

func (m *PetResponse) GetPet() *Pet {
	if m != nil {
		return m.Pet
	}
	return nil
}

type UserResponse struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserResponse) Reset()         { *m = UserResponse{} }
func (m *UserResponse) String() string { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()    {}
func (*UserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{5}
}
func (m *UserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserResponse.Unmarshal(m, b)
}
func (m *UserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserResponse.Marshal(b, m, deterministic)
}
func (dst *UserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserResponse.Merge(dst, src)
}
func (m *UserResponse) XXX_Size() int {
	return xxx_messageInfo_UserResponse.Size(m)
}
func (m *UserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UserResponse proto.InternalMessageInfo

func (m *UserResponse) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type PetRequest struct {
	Pet                  *Pet     `protobuf:"bytes,1,opt,name=pet,proto3" json:"pet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PetRequest) Reset()         { *m = PetRequest{} }
func (m *PetRequest) String() string { return proto.CompactTextString(m) }
func (*PetRequest) ProtoMessage()    {}
func (*PetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{6}
}
func (m *PetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PetRequest.Unmarshal(m, b)
}
func (m *PetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PetRequest.Marshal(b, m, deterministic)
}
func (dst *PetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PetRequest.Merge(dst, src)
}
func (m *PetRequest) XXX_Size() int {
	return xxx_messageInfo_PetRequest.Size(m)
}
func (m *PetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PetRequest proto.InternalMessageInfo

func (m *PetRequest) GetPet() *Pet {
	if m != nil {
		return m.Pet
	}
	return nil
}

type UserRequest struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserRequest) Reset()         { *m = UserRequest{} }
func (m *UserRequest) String() string { return proto.CompactTextString(m) }
func (*UserRequest) ProtoMessage()    {}
func (*UserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_petstore_f48f4004d62baad2, []int{7}
}
func (m *UserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserRequest.Unmarshal(m, b)
}
func (m *UserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserRequest.Marshal(b, m, deterministic)
}
func (dst *UserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserRequest.Merge(dst, src)
}
func (m *UserRequest) XXX_Size() int {
	return xxx_messageInfo_UserRequest.Size(m)
}
func (m *UserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserRequest proto.InternalMessageInfo

func (m *UserRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func init() {
	proto.RegisterType((*Pet)(nil), "rest2grpc.Pet")
	proto.RegisterType((*User)(nil), "rest2grpc.User")
	proto.RegisterType((*PetByIdRequest)(nil), "rest2grpc.PetByIdRequest")
	proto.RegisterType((*UserByNameRequest)(nil), "rest2grpc.UserByNameRequest")
	proto.RegisterType((*PetResponse)(nil), "rest2grpc.PetResponse")
	proto.RegisterType((*UserResponse)(nil), "rest2grpc.UserResponse")
	proto.RegisterType((*PetRequest)(nil), "rest2grpc.PetRequest")
	proto.RegisterType((*UserRequest)(nil), "rest2grpc.UserRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Rest2GRPCPetStoreServiceClient is the client API for Rest2GRPCPetStoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type Rest2GRPCPetStoreServiceClient interface {
	PetById(ctx context.Context, in *PetByIdRequest, opts ...grpc.CallOption) (*PetResponse, error)
	UserByName(ctx context.Context, in *UserByNameRequest, opts ...grpc.CallOption) (*UserResponse, error)
	PetPUT(ctx context.Context, in *PetRequest, opts ...grpc.CallOption) (*PetResponse, error)
	UserPUT(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type rest2GRPCPetStoreServiceClient struct {
	cc *grpc.ClientConn
}

func NewRest2GRPCPetStoreServiceClient(cc *grpc.ClientConn) Rest2GRPCPetStoreServiceClient {
	return &rest2GRPCPetStoreServiceClient{cc}
}

func (c *rest2GRPCPetStoreServiceClient) PetById(ctx context.Context, in *PetByIdRequest, opts ...grpc.CallOption) (*PetResponse, error) {
	out := new(PetResponse)
	err := c.cc.Invoke(ctx, "/rest2grpc.Rest2GRPCPetStoreService/PetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rest2GRPCPetStoreServiceClient) UserByName(ctx context.Context, in *UserByNameRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/rest2grpc.Rest2GRPCPetStoreService/UserByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rest2GRPCPetStoreServiceClient) PetPUT(ctx context.Context, in *PetRequest, opts ...grpc.CallOption) (*PetResponse, error) {
	out := new(PetResponse)
	err := c.cc.Invoke(ctx, "/rest2grpc.Rest2GRPCPetStoreService/PetPUT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rest2GRPCPetStoreServiceClient) UserPUT(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/rest2grpc.Rest2GRPCPetStoreService/UserPUT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Rest2GRPCPetStoreServiceServer is the server API for Rest2GRPCPetStoreService service.
type Rest2GRPCPetStoreServiceServer interface {
	PetById(context.Context, *PetByIdRequest) (*PetResponse, error)
	UserByName(context.Context, *UserByNameRequest) (*UserResponse, error)
	PetPUT(context.Context, *PetRequest) (*PetResponse, error)
	UserPUT(context.Context, *UserRequest) (*UserResponse, error)
}

func RegisterRest2GRPCPetStoreServiceServer(s *grpc.Server, srv Rest2GRPCPetStoreServiceServer) {
	s.RegisterService(&_Rest2GRPCPetStoreService_serviceDesc, srv)
}

func _Rest2GRPCPetStoreService_PetById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Rest2GRPCPetStoreServiceServer).PetById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rest2grpc.Rest2GRPCPetStoreService/PetById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Rest2GRPCPetStoreServiceServer).PetById(ctx, req.(*PetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rest2GRPCPetStoreService_UserByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Rest2GRPCPetStoreServiceServer).UserByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rest2grpc.Rest2GRPCPetStoreService/UserByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Rest2GRPCPetStoreServiceServer).UserByName(ctx, req.(*UserByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rest2GRPCPetStoreService_PetPUT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Rest2GRPCPetStoreServiceServer).PetPUT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rest2grpc.Rest2GRPCPetStoreService/PetPUT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Rest2GRPCPetStoreServiceServer).PetPUT(ctx, req.(*PetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rest2GRPCPetStoreService_UserPUT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Rest2GRPCPetStoreServiceServer).UserPUT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rest2grpc.Rest2GRPCPetStoreService/UserPUT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Rest2GRPCPetStoreServiceServer).UserPUT(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rest2GRPCPetStoreService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rest2grpc.Rest2GRPCPetStoreService",
	HandlerType: (*Rest2GRPCPetStoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PetById",
			Handler:    _Rest2GRPCPetStoreService_PetById_Handler,
		},
		{
			MethodName: "UserByName",
			Handler:    _Rest2GRPCPetStoreService_UserByName_Handler,
		},
		{
			MethodName: "PetPUT",
			Handler:    _Rest2GRPCPetStoreService_PetPUT_Handler,
		},
		{
			MethodName: "UserPUT",
			Handler:    _Rest2GRPCPetStoreService_UserPUT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "petstore.proto",
}

func init() { proto.RegisterFile("petstore.proto", fileDescriptor_petstore_f48f4004d62baad2) }

var fileDescriptor_petstore_f48f4004d62baad2 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4f, 0xc2, 0x30,
	0x18, 0xc6, 0xb3, 0x31, 0x40, 0x5e, 0xcc, 0x8c, 0x8d, 0xe2, 0x24, 0x1e, 0x96, 0x7a, 0xc1, 0xcb,
	0x48, 0xc6, 0xc1, 0xc4, 0x78, 0x82, 0x83, 0xf1, 0x62, 0x9a, 0x22, 0x57, 0x13, 0x84, 0x37, 0xba,
	0x44, 0x58, 0x6d, 0x8b, 0x09, 0x5f, 0xc3, 0x4f, 0x6c, 0xda, 0xc1, 0xdc, 0x9f, 0x04, 0xbd, 0xad,
	0x7d, 0x9f, 0xe7, 0xfd, 0x3d, 0x7b, 0x36, 0xf0, 0x05, 0x6a, 0xa5, 0x53, 0x89, 0x91, 0x90, 0xa9,
	0x4e, 0x49, 0x47, 0xa2, 0xd2, 0xf1, 0x9b, 0x14, 0x0b, 0x7a, 0x03, 0x0d, 0x86, 0x9a, 0xf8, 0xe0,
	0x26, 0xcb, 0xc0, 0x09, 0x9d, 0x41, 0x93, 0xbb, 0xc9, 0x92, 0x10, 0xf0, 0xd6, 0xf3, 0x15, 0x06,
	0x6e, 0xe8, 0x0c, 0x3a, 0xdc, 0x3e, 0xd3, 0x17, 0xf0, 0x66, 0x0a, 0x65, 0x4d, 0xdb, 0x87, 0xa3,
	0x8d, 0x42, 0x59, 0xd0, 0xe7, 0x67, 0x72, 0x06, 0x4d, 0x5c, 0xcd, 0x93, 0x8f, 0xa0, 0x61, 0x07,
	0xd9, 0xc1, 0xdc, 0x8a, 0xf7, 0x74, 0x8d, 0x81, 0x97, 0xdd, 0xda, 0x03, 0x0d, 0xc1, 0x67, 0xa8,
	0xc7, 0xdb, 0xc7, 0x25, 0xc7, 0xcf, 0x0d, 0xaa, 0x5a, 0x2a, 0x3a, 0x84, 0x53, 0x93, 0x60, 0xbc,
	0x7d, 0x9a, 0xaf, 0x70, 0x2f, 0x2a, 0xe2, 0x9d, 0x32, 0x9e, 0x0e, 0xa1, 0xcb, 0x50, 0x73, 0x54,
	0x22, 0x5d, 0x2b, 0x24, 0x21, 0x34, 0x04, 0x6a, 0xab, 0xea, 0xc6, 0x7e, 0x94, 0xb7, 0x10, 0x19,
	0x91, 0x19, 0xd1, 0x11, 0x1c, 0x1b, 0x42, 0xee, 0xb8, 0x06, 0xcf, 0x2c, 0xdb, 0x59, 0x4e, 0x0a,
	0x16, 0x2b, 0xb3, 0x43, 0x1a, 0x01, 0x58, 0x4a, 0x96, 0xe7, 0x6f, 0x48, 0x0c, 0xdd, 0x0c, 0x92,
	0x19, 0xfe, 0xc3, 0x88, 0xbf, 0x5d, 0x08, 0xb8, 0x19, 0x3c, 0x70, 0x36, 0x61, 0xa8, 0xa7, 0xe6,
	0x73, 0x4e, 0x51, 0x7e, 0x25, 0x0b, 0x24, 0xf7, 0xd0, 0xde, 0x35, 0x47, 0x2e, 0xcb, 0xc0, 0x42,
	0x9b, 0xfd, 0x5e, 0x25, 0xcb, 0xfe, 0x1d, 0x27, 0x00, 0xbf, 0xad, 0x92, 0xab, 0x0a, 0xbf, 0x54,
	0x76, 0xff, 0xa2, 0x9a, 0x6e, 0xbf, 0xe4, 0x16, 0x5a, 0x0c, 0x35, 0x9b, 0x3d, 0x93, 0xf3, 0x2a,
	0xe6, 0x30, 0xfd, 0x0e, 0xda, 0x66, 0x91, 0x71, 0xf6, 0x6a, 0xcb, 0x0f, 0x43, 0x5f, 0x5b, 0xf6,
	0x77, 0x1e, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x09, 0x68, 0x48, 0xe2, 0xe0, 0x02, 0x00, 0x00,
}
