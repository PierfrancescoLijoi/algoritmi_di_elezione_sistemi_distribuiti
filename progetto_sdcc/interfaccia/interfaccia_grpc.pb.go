// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: interfaccia.proto

package interfaccia

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	RingElection_SendElection_FullMethodName = "/RingElection/SendElection"
	RingElection_SendResult_FullMethodName   = "/RingElection/SendResult"
	RingElection_SendCheck_FullMethodName    = "/RingElection/SendCheck"
	RingElection_SendRequest_FullMethodName  = "/RingElection/SendRequest"
	RingElection_SendUpdate_FullMethodName   = "/RingElection/SendUpdate"
)

// RingElectionClient is the client API for RingElection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RingElectionClient interface {
	SendElection(ctx context.Context, in *MessageElection, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendResult(ctx context.Context, in *ResultElection, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendCheck(ctx context.Context, in *CheckLeader, opts ...grpc.CallOption) (*CheckLeader, error)
	SendRequest(ctx context.Context, in *JSONClientRegistration, opts ...grpc.CallOption) (*JSONListPeer, error)
	SendUpdate(ctx context.Context, in *JSONListPeer, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type ringElectionClient struct {
	cc grpc.ClientConnInterface
}

func NewRingElectionClient(cc grpc.ClientConnInterface) RingElectionClient {
	return &ringElectionClient{cc}
}

func (c *ringElectionClient) SendElection(ctx context.Context, in *MessageElection, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RingElection_SendElection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringElectionClient) SendResult(ctx context.Context, in *ResultElection, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RingElection_SendResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringElectionClient) SendCheck(ctx context.Context, in *CheckLeader, opts ...grpc.CallOption) (*CheckLeader, error) {
	out := new(CheckLeader)
	err := c.cc.Invoke(ctx, RingElection_SendCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringElectionClient) SendRequest(ctx context.Context, in *JSONClientRegistration, opts ...grpc.CallOption) (*JSONListPeer, error) {
	out := new(JSONListPeer)
	err := c.cc.Invoke(ctx, RingElection_SendRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ringElectionClient) SendUpdate(ctx context.Context, in *JSONListPeer, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RingElection_SendUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RingElectionServer is the server API for RingElection service.
// All implementations must embed UnimplementedRingElectionServer
// for forward compatibility
type RingElectionServer interface {
	SendElection(context.Context, *MessageElection) (*emptypb.Empty, error)
	SendResult(context.Context, *ResultElection) (*emptypb.Empty, error)
	SendCheck(context.Context, *CheckLeader) (*CheckLeader, error)
	SendRequest(context.Context, *JSONClientRegistration) (*JSONListPeer, error)
	SendUpdate(context.Context, *JSONListPeer) (*emptypb.Empty, error)
	mustEmbedUnimplementedRingElectionServer()
}

// UnimplementedRingElectionServer must be embedded to have forward compatible implementations.
type UnimplementedRingElectionServer struct {
}

func (UnimplementedRingElectionServer) SendElection(context.Context, *MessageElection) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendElection not implemented")
}
func (UnimplementedRingElectionServer) SendResult(context.Context, *ResultElection) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendResult not implemented")
}
func (UnimplementedRingElectionServer) SendCheck(context.Context, *CheckLeader) (*CheckLeader, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCheck not implemented")
}
func (UnimplementedRingElectionServer) SendRequest(context.Context, *JSONClientRegistration) (*JSONListPeer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRequest not implemented")
}
func (UnimplementedRingElectionServer) SendUpdate(context.Context, *JSONListPeer) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUpdate not implemented")
}
func (UnimplementedRingElectionServer) mustEmbedUnimplementedRingElectionServer() {}

// UnsafeRingElectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RingElectionServer will
// result in compilation errors.
type UnsafeRingElectionServer interface {
	mustEmbedUnimplementedRingElectionServer()
}

func RegisterRingElectionServer(s grpc.ServiceRegistrar, srv RingElectionServer) {
	s.RegisterService(&RingElection_ServiceDesc, srv)
}

func _RingElection_SendElection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageElection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingElectionServer).SendElection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingElection_SendElection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingElectionServer).SendElection(ctx, req.(*MessageElection))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingElection_SendResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultElection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingElectionServer).SendResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingElection_SendResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingElectionServer).SendResult(ctx, req.(*ResultElection))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingElection_SendCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckLeader)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingElectionServer).SendCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingElection_SendCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingElectionServer).SendCheck(ctx, req.(*CheckLeader))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingElection_SendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONClientRegistration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingElectionServer).SendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingElection_SendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingElectionServer).SendRequest(ctx, req.(*JSONClientRegistration))
	}
	return interceptor(ctx, in, info, handler)
}

func _RingElection_SendUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONListPeer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RingElectionServer).SendUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RingElection_SendUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RingElectionServer).SendUpdate(ctx, req.(*JSONListPeer))
	}
	return interceptor(ctx, in, info, handler)
}

// RingElection_ServiceDesc is the grpc.ServiceDesc for RingElection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RingElection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RingElection",
	HandlerType: (*RingElectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendElection",
			Handler:    _RingElection_SendElection_Handler,
		},
		{
			MethodName: "SendResult",
			Handler:    _RingElection_SendResult_Handler,
		},
		{
			MethodName: "SendCheck",
			Handler:    _RingElection_SendCheck_Handler,
		},
		{
			MethodName: "SendRequest",
			Handler:    _RingElection_SendRequest_Handler,
		},
		{
			MethodName: "SendUpdate",
			Handler:    _RingElection_SendUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interfaccia.proto",
}

const (
	RaftElection_SendRequest_FullMethodName      = "/RaftElection/SendRequest"
	RaftElection_SendUpdate_FullMethodName       = "/RaftElection/SendUpdate"
	RaftElection_SendRaftElection_FullMethodName = "/RaftElection/SendRaftElection"
	RaftElection_SendHeartBeat_FullMethodName    = "/RaftElection/sendHeartBeat"
)

// RaftElectionClient is the client API for RaftElection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RaftElectionClient interface {
	SendRequest(ctx context.Context, in *JSONClientRegistration, opts ...grpc.CallOption) (*JSONListPeer, error)
	SendUpdate(ctx context.Context, in *JSONListPeer, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendRaftElection(ctx context.Context, in *RequestElection, opts ...grpc.CallOption) (*AnswerElection, error)
	SendHeartBeat(ctx context.Context, in *HeartBeat, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type raftElectionClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftElectionClient(cc grpc.ClientConnInterface) RaftElectionClient {
	return &raftElectionClient{cc}
}

func (c *raftElectionClient) SendRequest(ctx context.Context, in *JSONClientRegistration, opts ...grpc.CallOption) (*JSONListPeer, error) {
	out := new(JSONListPeer)
	err := c.cc.Invoke(ctx, RaftElection_SendRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftElectionClient) SendUpdate(ctx context.Context, in *JSONListPeer, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RaftElection_SendUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftElectionClient) SendRaftElection(ctx context.Context, in *RequestElection, opts ...grpc.CallOption) (*AnswerElection, error) {
	out := new(AnswerElection)
	err := c.cc.Invoke(ctx, RaftElection_SendRaftElection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftElectionClient) SendHeartBeat(ctx context.Context, in *HeartBeat, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RaftElection_SendHeartBeat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftElectionServer is the server API for RaftElection service.
// All implementations must embed UnimplementedRaftElectionServer
// for forward compatibility
type RaftElectionServer interface {
	SendRequest(context.Context, *JSONClientRegistration) (*JSONListPeer, error)
	SendUpdate(context.Context, *JSONListPeer) (*emptypb.Empty, error)
	SendRaftElection(context.Context, *RequestElection) (*AnswerElection, error)
	SendHeartBeat(context.Context, *HeartBeat) (*emptypb.Empty, error)
	mustEmbedUnimplementedRaftElectionServer()
}

// UnimplementedRaftElectionServer must be embedded to have forward compatible implementations.
type UnimplementedRaftElectionServer struct {
}

func (UnimplementedRaftElectionServer) SendRequest(context.Context, *JSONClientRegistration) (*JSONListPeer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRequest not implemented")
}
func (UnimplementedRaftElectionServer) SendUpdate(context.Context, *JSONListPeer) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUpdate not implemented")
}
func (UnimplementedRaftElectionServer) SendRaftElection(context.Context, *RequestElection) (*AnswerElection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRaftElection not implemented")
}
func (UnimplementedRaftElectionServer) SendHeartBeat(context.Context, *HeartBeat) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHeartBeat not implemented")
}
func (UnimplementedRaftElectionServer) mustEmbedUnimplementedRaftElectionServer() {}

// UnsafeRaftElectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftElectionServer will
// result in compilation errors.
type UnsafeRaftElectionServer interface {
	mustEmbedUnimplementedRaftElectionServer()
}

func RegisterRaftElectionServer(s grpc.ServiceRegistrar, srv RaftElectionServer) {
	s.RegisterService(&RaftElection_ServiceDesc, srv)
}

func _RaftElection_SendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONClientRegistration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftElectionServer).SendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftElection_SendRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftElectionServer).SendRequest(ctx, req.(*JSONClientRegistration))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftElection_SendUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JSONListPeer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftElectionServer).SendUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftElection_SendUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftElectionServer).SendUpdate(ctx, req.(*JSONListPeer))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftElection_SendRaftElection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestElection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftElectionServer).SendRaftElection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftElection_SendRaftElection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftElectionServer).SendRaftElection(ctx, req.(*RequestElection))
	}
	return interceptor(ctx, in, info, handler)
}

func _RaftElection_SendHeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftElectionServer).SendHeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RaftElection_SendHeartBeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftElectionServer).SendHeartBeat(ctx, req.(*HeartBeat))
	}
	return interceptor(ctx, in, info, handler)
}

// RaftElection_ServiceDesc is the grpc.ServiceDesc for RaftElection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RaftElection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RaftElection",
	HandlerType: (*RaftElectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendRequest",
			Handler:    _RaftElection_SendRequest_Handler,
		},
		{
			MethodName: "SendUpdate",
			Handler:    _RaftElection_SendUpdate_Handler,
		},
		{
			MethodName: "SendRaftElection",
			Handler:    _RaftElection_SendRaftElection_Handler,
		},
		{
			MethodName: "sendHeartBeat",
			Handler:    _RaftElection_SendHeartBeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interfaccia.proto",
}
