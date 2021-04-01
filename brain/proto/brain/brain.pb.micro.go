// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/brain/brain.proto

package micro_chatting_service_brain

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Brain service

func NewBrainEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Brain service

type BrainService interface {
	Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Brain_StreamService, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Brain_PingPongService, error)
}

type brainService struct {
	c    client.Client
	name string
}

func NewBrainService(name string, c client.Client) BrainService {
	return &brainService{
		c:    c,
		name: name,
	}
}

func (c *brainService) Call(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Brain.Call", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *brainService) Stream(ctx context.Context, in *StreamingRequest, opts ...client.CallOption) (Brain_StreamService, error) {
	req := c.c.NewRequest(c.name, "Brain.Stream", &StreamingRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &brainServiceStream{stream}, nil
}

type Brain_StreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Recv() (*StreamingResponse, error)
}

type brainServiceStream struct {
	stream client.Stream
}

func (x *brainServiceStream) Close() error {
	return x.stream.Close()
}

func (x *brainServiceStream) Context() context.Context {
	return x.stream.Context()
}

func (x *brainServiceStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brainServiceStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *brainServiceStream) Recv() (*StreamingResponse, error) {
	m := new(StreamingResponse)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *brainService) PingPong(ctx context.Context, opts ...client.CallOption) (Brain_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Brain.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &brainServicePingPong{stream}, nil
}

type Brain_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type brainServicePingPong struct {
	stream client.Stream
}

func (x *brainServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *brainServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *brainServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brainServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *brainServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *brainServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Brain service

type BrainHandler interface {
	Call(context.Context, *Request, *Response) error
	Stream(context.Context, *StreamingRequest, Brain_StreamStream) error
	PingPong(context.Context, Brain_PingPongStream) error
}

func RegisterBrainHandler(s server.Server, hdlr BrainHandler, opts ...server.HandlerOption) error {
	type brain interface {
		Call(ctx context.Context, in *Request, out *Response) error
		Stream(ctx context.Context, stream server.Stream) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Brain struct {
		brain
	}
	h := &brainHandler{hdlr}
	return s.Handle(s.NewHandler(&Brain{h}, opts...))
}

type brainHandler struct {
	BrainHandler
}

func (h *brainHandler) Call(ctx context.Context, in *Request, out *Response) error {
	return h.BrainHandler.Call(ctx, in, out)
}

func (h *brainHandler) Stream(ctx context.Context, stream server.Stream) error {
	m := new(StreamingRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.BrainHandler.Stream(ctx, m, &brainStreamStream{stream})
}

type Brain_StreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*StreamingResponse) error
}

type brainStreamStream struct {
	stream server.Stream
}

func (x *brainStreamStream) Close() error {
	return x.stream.Close()
}

func (x *brainStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *brainStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brainStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *brainStreamStream) Send(m *StreamingResponse) error {
	return x.stream.Send(m)
}

func (h *brainHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.BrainHandler.PingPong(ctx, &brainPingPongStream{stream})
}

type Brain_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type brainPingPongStream struct {
	stream server.Stream
}

func (x *brainPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *brainPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *brainPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *brainPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *brainPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *brainPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}
