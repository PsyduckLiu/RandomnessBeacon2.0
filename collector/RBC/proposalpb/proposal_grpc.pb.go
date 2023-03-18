// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proposal.proto

package proposalpb

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

// ProposalHandleClient is the client API for ProposalHandle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProposalHandleClient interface {
	// Handles a received proposal message
	ProposalReceive(ctx context.Context, in *Proposal, opts ...grpc.CallOption) (*ProposalResponse, error)
}

type proposalHandleClient struct {
	cc grpc.ClientConnInterface
}

func NewProposalHandleClient(cc grpc.ClientConnInterface) ProposalHandleClient {
	return &proposalHandleClient{cc}
}

func (c *proposalHandleClient) ProposalReceive(ctx context.Context, in *Proposal, opts ...grpc.CallOption) (*ProposalResponse, error) {
	out := new(ProposalResponse)
	err := c.cc.Invoke(ctx, "/proposalpb.ProposalHandle/ProposalReceive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProposalHandleServer is the server API for ProposalHandle service.
// All implementations must embed UnimplementedProposalHandleServer
// for forward compatibility
type ProposalHandleServer interface {
	// Handles a received proposal message
	ProposalReceive(context.Context, *Proposal) (*ProposalResponse, error)
	mustEmbedUnimplementedProposalHandleServer()
}

// UnimplementedProposalHandleServer must be embedded to have forward compatible implementations.
type UnimplementedProposalHandleServer struct {
}

func (UnimplementedProposalHandleServer) ProposalReceive(context.Context, *Proposal) (*ProposalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProposalReceive not implemented")
}
func (UnimplementedProposalHandleServer) mustEmbedUnimplementedProposalHandleServer() {}

// UnsafeProposalHandleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProposalHandleServer will
// result in compilation errors.
type UnsafeProposalHandleServer interface {
	mustEmbedUnimplementedProposalHandleServer()
}

func RegisterProposalHandleServer(s grpc.ServiceRegistrar, srv ProposalHandleServer) {
	s.RegisterService(&ProposalHandle_ServiceDesc, srv)
}

func _ProposalHandle_ProposalReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Proposal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProposalHandleServer).ProposalReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proposalpb.ProposalHandle/ProposalReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProposalHandleServer).ProposalReceive(ctx, req.(*Proposal))
	}
	return interceptor(ctx, in, info, handler)
}

// ProposalHandle_ServiceDesc is the grpc.ServiceDesc for ProposalHandle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProposalHandle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proposalpb.ProposalHandle",
	HandlerType: (*ProposalHandleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProposalReceive",
			Handler:    _ProposalHandle_ProposalReceive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proposal.proto",
}
