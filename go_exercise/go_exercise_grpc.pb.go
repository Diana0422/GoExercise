// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package go_exercise

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

// GrepMapReduceClient is the client API for GrepMapReduce service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrepMapReduceClient interface {
	// Splits the file and assigns each chunk to
	//a different worker to process it. Then collects
	//the results and returns the grep file
	Grep(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error)
}

type grepMapReduceClient struct {
	cc grpc.ClientConnInterface
}

func NewGrepMapReduceClient(cc grpc.ClientConnInterface) GrepMapReduceClient {
	return &grepMapReduceClient{cc}
}

func (c *grepMapReduceClient) Grep(ctx context.Context, in *File, opts ...grpc.CallOption) (*File, error) {
	out := new(File)
	err := c.cc.Invoke(ctx, "/go_exercise.GrepMapReduce/Grep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrepMapReduceServer is the server API for GrepMapReduce service.
// All implementations must embed UnimplementedGrepMapReduceServer
// for forward compatibility
type GrepMapReduceServer interface {
	// Splits the file and assigns each chunk to
	//a different worker to process it. Then collects
	//the results and returns the grep file
	Grep(context.Context, *File) (*File, error)
	mustEmbedUnimplementedGrepMapReduceServer()
}

// UnimplementedGrepMapReduceServer must be embedded to have forward compatible implementations.
type UnimplementedGrepMapReduceServer struct {
}

func (UnimplementedGrepMapReduceServer) Grep(context.Context, *File) (*File, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Grep not implemented")
}
func (UnimplementedGrepMapReduceServer) mustEmbedUnimplementedGrepMapReduceServer() {}

// UnsafeGrepMapReduceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrepMapReduceServer will
// result in compilation errors.
type UnsafeGrepMapReduceServer interface {
	mustEmbedUnimplementedGrepMapReduceServer()
}

func RegisterGrepMapReduceServer(s grpc.ServiceRegistrar, srv GrepMapReduceServer) {
	s.RegisterService(&GrepMapReduce_ServiceDesc, srv)
}

func _GrepMapReduce_Grep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(File)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrepMapReduceServer).Grep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/go_exercise.GrepMapReduce/Grep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrepMapReduceServer).Grep(ctx, req.(*File))
	}
	return interceptor(ctx, in, info, handler)
}

// GrepMapReduce_ServiceDesc is the grpc.ServiceDesc for GrepMapReduce service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrepMapReduce_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "go_exercise.GrepMapReduce",
	HandlerType: (*GrepMapReduceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Grep",
			Handler:    _GrepMapReduce_Grep_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "go_exercise/go_exercise.proto",
}