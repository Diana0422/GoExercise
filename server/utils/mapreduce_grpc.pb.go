// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package utils

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

// MapReduceClient is the client API for MapReduce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapReduceClient interface {
	Grep(ctx context.Context, in *FileChunk, opts ...grpc.CallOption) (MapReduce_GrepClient, error)
}

type mapReduceClient struct {
	cc grpc.ClientConnInterface
}

func NewMapReduceClient(cc grpc.ClientConnInterface) MapReduceClient {
	return &mapReduceClient{cc}
}

func (c *mapReduceClient) Grep(ctx context.Context, in *FileChunk, opts ...grpc.CallOption) (MapReduce_GrepClient, error) {
	stream, err := c.cc.NewStream(ctx, &MapReduce_ServiceDesc.Streams[0], "/utils.MapReduce/grep", opts...)
	if err != nil {
		return nil, err
	}
	x := &mapReduceGrepClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MapReduce_GrepClient interface {
	Recv() (*GrepRow, error)
	grpc.ClientStream
}

type mapReduceGrepClient struct {
	grpc.ClientStream
}

func (x *mapReduceGrepClient) Recv() (*GrepRow, error) {
	m := new(GrepRow)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MapReduceServer is the server API for MapReduce service.
// All implementations must embed UnimplementedMapReduceServer
// for forward compatibility
type MapReduceServer interface {
	Grep(*FileChunk, MapReduce_GrepServer) error
	mustEmbedUnimplementedMapReduceServer()
}

// UnimplementedMapReduceServer must be embedded to have forward compatible implementations.
type UnimplementedMapReduceServer struct {
}

func (UnimplementedMapReduceServer) Grep(*FileChunk, MapReduce_GrepServer) error {
	return status.Errorf(codes.Unimplemented, "method Grep not implemented")
}
func (UnimplementedMapReduceServer) mustEmbedUnimplementedMapReduceServer() {}

// UnsafeMapReduceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapReduceServer will
// result in compilation errors.
type UnsafeMapReduceServer interface {
	mustEmbedUnimplementedMapReduceServer()
}

func RegisterMapReduceServer(s grpc.ServiceRegistrar, srv MapReduceServer) {
	s.RegisterService(&MapReduce_ServiceDesc, srv)
}

func _MapReduce_Grep_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileChunk)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MapReduceServer).Grep(m, &mapReduceGrepServer{stream})
}

type MapReduce_GrepServer interface {
	Send(*GrepRow) error
	grpc.ServerStream
}

type mapReduceGrepServer struct {
	grpc.ServerStream
}

func (x *mapReduceGrepServer) Send(m *GrepRow) error {
	return x.ServerStream.SendMsg(m)
}

// MapReduce_ServiceDesc is the grpc.ServiceDesc for MapReduce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapReduce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "utils.MapReduce",
	HandlerType: (*MapReduceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "grep",
			Handler:       _MapReduce_Grep_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "utils/mapreduce.proto",
}