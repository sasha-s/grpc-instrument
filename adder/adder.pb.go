// Code generated by protoc-gen-go.
// source: adder/adder.proto
// DO NOT EDIT!

/*
Package adder is a generated protocol buffer package.

It is generated from these files:
	adder/adder.proto

It has these top-level messages:
	AddRequest
	AddReply
	Add2Request
	Add2Reply
*/
package adder

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

type AddRequest struct {
	A int32 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b" json:"b,omitempty"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}

type AddReply struct {
	R int32 `protobuf:"varint,1,opt,name=r" json:"r,omitempty"`
}

func (m *AddReply) Reset()         { *m = AddReply{} }
func (m *AddReply) String() string { return proto.CompactTextString(m) }
func (*AddReply) ProtoMessage()    {}

type Add2Request struct {
	A int32 `protobuf:"varint,1,opt,name=a" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b" json:"b,omitempty"`
}

func (m *Add2Request) Reset()         { *m = Add2Request{} }
func (m *Add2Request) String() string { return proto.CompactTextString(m) }
func (*Add2Request) ProtoMessage()    {}

type Add2Reply struct {
	R int32 `protobuf:"varint,1,opt,name=r" json:"r,omitempty"`
}

func (m *Add2Reply) Reset()         { *m = Add2Reply{} }
func (m *Add2Reply) String() string { return proto.CompactTextString(m) }
func (*Add2Reply) ProtoMessage()    {}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Adder service

type AdderClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error)
	Add2(ctx context.Context, in *Add2Request, opts ...grpc.CallOption) (*Add2Reply, error)
}

type adderClient struct {
	cc *grpc.ClientConn
}

func NewAdderClient(cc *grpc.ClientConn) AdderClient {
	return &adderClient{cc}
}

func (c *adderClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddReply, error) {
	out := new(AddReply)
	err := grpc.Invoke(ctx, "/adder.Adder/Add", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adderClient) Add2(ctx context.Context, in *Add2Request, opts ...grpc.CallOption) (*Add2Reply, error) {
	out := new(Add2Reply)
	err := grpc.Invoke(ctx, "/adder.Adder/Add2", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Adder service

type AdderServer interface {
	Add(context.Context, *AddRequest) (*AddReply, error)
	Add2(context.Context, *Add2Request) (*Add2Reply, error)
}

func RegisterAdderServer(s *grpc.Server, srv AdderServer) {
	s.RegisterService(&_Adder_serviceDesc, srv)
}

func _Adder_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(AdderServer).Add(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Adder_Add2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Add2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(AdderServer).Add2(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Adder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "adder.Adder",
	HandlerType: (*AdderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Adder_Add_Handler,
		},
		{
			MethodName: "Add2",
			Handler:    _Adder_Add2_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
