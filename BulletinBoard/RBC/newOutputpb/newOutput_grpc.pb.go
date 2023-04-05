// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: newOutput.proto

package newOutputpb

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

// NewOutputHandleClient is the client API for NewOutputHandle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NewOutputHandleClient interface {
	// Handles a received newOutput message
	NewOutputReceive(ctx context.Context, in *NewOutput, opts ...grpc.CallOption) (*NewOutputResponse, error)
}

type newOutputHandleClient struct {
	cc grpc.ClientConnInterface
}

func NewNewOutputHandleClient(cc grpc.ClientConnInterface) NewOutputHandleClient {
	return &newOutputHandleClient{cc}
}

func (c *newOutputHandleClient) NewOutputReceive(ctx context.Context, in *NewOutput, opts ...grpc.CallOption) (*NewOutputResponse, error) {
	out := new(NewOutputResponse)
	err := c.cc.Invoke(ctx, "/newOutputpb.NewOutputHandle/NewOutputReceive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NewOutputHandleServer is the server API for NewOutputHandle service.
// All implementations must embed UnimplementedNewOutputHandleServer
// for forward compatibility
type NewOutputHandleServer interface {
	// Handles a received newOutput message
	NewOutputReceive(context.Context, *NewOutput) (*NewOutputResponse, error)
	mustEmbedUnimplementedNewOutputHandleServer()
}

// UnimplementedNewOutputHandleServer must be embedded to have forward compatible implementations.
type UnimplementedNewOutputHandleServer struct {
}

func (UnimplementedNewOutputHandleServer) NewOutputReceive(context.Context, *NewOutput) (*NewOutputResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewOutputReceive not implemented")
}
func (UnimplementedNewOutputHandleServer) mustEmbedUnimplementedNewOutputHandleServer() {}

// UnsafeNewOutputHandleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NewOutputHandleServer will
// result in compilation errors.
type UnsafeNewOutputHandleServer interface {
	mustEmbedUnimplementedNewOutputHandleServer()
}

func RegisterNewOutputHandleServer(s grpc.ServiceRegistrar, srv NewOutputHandleServer) {
	s.RegisterService(&NewOutputHandle_ServiceDesc, srv)
}

func _NewOutputHandle_NewOutputReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewOutput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NewOutputHandleServer).NewOutputReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/newOutputpb.NewOutputHandle/NewOutputReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NewOutputHandleServer).NewOutputReceive(ctx, req.(*NewOutput))
	}
	return interceptor(ctx, in, info, handler)
}

// NewOutputHandle_ServiceDesc is the grpc.ServiceDesc for NewOutputHandle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NewOutputHandle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "newOutputpb.NewOutputHandle",
	HandlerType: (*NewOutputHandleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewOutputReceive",
			Handler:    _NewOutputHandle_NewOutputReceive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "newOutput.proto",
}