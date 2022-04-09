// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package message

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MessageClient is the client API for Message service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageClient interface {
	MessagePeer(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	SubscribeToPeer(ctx context.Context, opts ...grpc.CallOption) (Message_SubscribeToPeerClient, error)
}

type messageClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageClient(cc grpc.ClientConnInterface) MessageClient {
	return &messageClient{cc}
}

func (c *messageClient) MessagePeer(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/Message/MessagePeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageClient) SubscribeToPeer(ctx context.Context, opts ...grpc.CallOption) (Message_SubscribeToPeerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Message_ServiceDesc.Streams[0], "/Message/SubscribeToPeer", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageSubscribeToPeerClient{stream}
	return x, nil
}

type Message_SubscribeToPeerClient interface {
	Send(*MessageRequest) error
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type messageSubscribeToPeerClient struct {
	grpc.ClientStream
}

func (x *messageSubscribeToPeerClient) Send(m *MessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *messageSubscribeToPeerClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServer is the server API for Message service.
// All implementations must embed UnimplementedMessageServer
// for forward compatibility
type MessageServer interface {
	MessagePeer(context.Context, *MessageRequest) (*MessageResponse, error)
	SubscribeToPeer(Message_SubscribeToPeerServer) error
	mustEmbedUnimplementedMessageServer()
}

// UnimplementedMessageServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServer struct {
}

func (UnimplementedMessageServer) MessagePeer(context.Context, *MessageRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessagePeer not implemented")
}
func (UnimplementedMessageServer) SubscribeToPeer(Message_SubscribeToPeerServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToPeer not implemented")
}
func (UnimplementedMessageServer) mustEmbedUnimplementedMessageServer() {}

// UnsafeMessageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServer will
// result in compilation errors.
type UnsafeMessageServer interface {
	mustEmbedUnimplementedMessageServer()
}

func RegisterMessageServer(s grpc.ServiceRegistrar, srv MessageServer) {
	s.RegisterService(&Message_ServiceDesc, srv)
}

func _Message_MessagePeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServer).MessagePeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Message/MessagePeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServer).MessagePeer(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Message_SubscribeToPeer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MessageServer).SubscribeToPeer(&messageSubscribeToPeerServer{stream})
}

type Message_SubscribeToPeerServer interface {
	Send(*MessageResponse) error
	Recv() (*MessageRequest, error)
	grpc.ServerStream
}

type messageSubscribeToPeerServer struct {
	grpc.ServerStream
}

func (x *messageSubscribeToPeerServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *messageSubscribeToPeerServer) Recv() (*MessageRequest, error) {
	m := new(MessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Message_ServiceDesc is the grpc.ServiceDesc for Message service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Message_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Message",
	HandlerType: (*MessageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MessagePeer",
			Handler:    _Message_MessagePeer_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToPeer",
			Handler:       _Message_SubscribeToPeer_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "message.proto",
}
