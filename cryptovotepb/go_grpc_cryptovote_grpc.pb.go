// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package cryptovotepb

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

// CryptoVoteServiceClient is the client API for CryptoVoteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CryptoVoteServiceClient interface {
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	CreateCryptoVote(ctx context.Context, in *CreateCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error)
	// Recupera dados de uma CryptoVote usando Server-Side Streaming RPC.
	RetrieveAllCryptoVoteByFilter(ctx context.Context, in *RetrieveCryptoReq, opts ...grpc.CallOption) (CryptoVoteService_RetrieveAllCryptoVoteByFilterClient, error)
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	UpdateOneCryptoVoteByFilter(ctx context.Context, in *UpdateCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error)
	// Atualiza dados de Qtd_Upvote de uma CryptoVote usando um simple RPC.
	AddUpVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*TotalChangesRes, error)
	// Atualiza dados de Qtd_Upvote de uma CryptoVote usando um simple RPC.
	AddDownVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*TotalChangesRes, error)
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	DeleteAllCryptoVoteByFilter(ctx context.Context, in *DeleteCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error)
}

type cryptoVoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCryptoVoteServiceClient(cc grpc.ClientConnInterface) CryptoVoteServiceClient {
	return &cryptoVoteServiceClient{cc}
}

func (c *cryptoVoteServiceClient) CreateCryptoVote(ctx context.Context, in *CreateCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error) {
	out := new(TotalChangesRes)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/CreateCryptoVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) RetrieveAllCryptoVoteByFilter(ctx context.Context, in *RetrieveCryptoReq, opts ...grpc.CallOption) (CryptoVoteService_RetrieveAllCryptoVoteByFilterClient, error) {
	stream, err := c.cc.NewStream(ctx, &CryptoVoteService_ServiceDesc.Streams[0], "/cryptovotepb.CryptoVoteService/RetrieveAllCryptoVoteByFilter", opts...)
	if err != nil {
		return nil, err
	}
	x := &cryptoVoteServiceRetrieveAllCryptoVoteByFilterClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CryptoVoteService_RetrieveAllCryptoVoteByFilterClient interface {
	Recv() (*CryptoVote, error)
	grpc.ClientStream
}

type cryptoVoteServiceRetrieveAllCryptoVoteByFilterClient struct {
	grpc.ClientStream
}

func (x *cryptoVoteServiceRetrieveAllCryptoVoteByFilterClient) Recv() (*CryptoVote, error) {
	m := new(CryptoVote)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cryptoVoteServiceClient) UpdateOneCryptoVoteByFilter(ctx context.Context, in *UpdateCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error) {
	out := new(TotalChangesRes)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/UpdateOneCryptoVoteByFilter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) AddUpVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*TotalChangesRes, error) {
	out := new(TotalChangesRes)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/AddUpVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) AddDownVote(ctx context.Context, in *AddVoteReq, opts ...grpc.CallOption) (*TotalChangesRes, error) {
	out := new(TotalChangesRes)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/AddDownVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) DeleteAllCryptoVoteByFilter(ctx context.Context, in *DeleteCryptoReq, opts ...grpc.CallOption) (*TotalChangesRes, error) {
	out := new(TotalChangesRes)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/DeleteAllCryptoVoteByFilter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CryptoVoteServiceServer is the server API for CryptoVoteService service.
// All implementations must embed UnimplementedCryptoVoteServiceServer
// for forward compatibility
type CryptoVoteServiceServer interface {
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	CreateCryptoVote(context.Context, *CreateCryptoReq) (*TotalChangesRes, error)
	// Recupera dados de uma CryptoVote usando Server-Side Streaming RPC.
	RetrieveAllCryptoVoteByFilter(*RetrieveCryptoReq, CryptoVoteService_RetrieveAllCryptoVoteByFilterServer) error
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	UpdateOneCryptoVoteByFilter(context.Context, *UpdateCryptoReq) (*TotalChangesRes, error)
	// Atualiza dados de Qtd_Upvote de uma CryptoVote usando um simple RPC.
	AddUpVote(context.Context, *AddVoteReq) (*TotalChangesRes, error)
	// Atualiza dados de Qtd_Upvote de uma CryptoVote usando um simple RPC.
	AddDownVote(context.Context, *AddVoteReq) (*TotalChangesRes, error)
	// Atualiza dados de uma CryptoVote usando um simple RPC.
	DeleteAllCryptoVoteByFilter(context.Context, *DeleteCryptoReq) (*TotalChangesRes, error)
	mustEmbedUnimplementedCryptoVoteServiceServer()
}

// UnimplementedCryptoVoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCryptoVoteServiceServer struct {
}

func (UnimplementedCryptoVoteServiceServer) CreateCryptoVote(context.Context, *CreateCryptoReq) (*TotalChangesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCryptoVote not implemented")
}
func (UnimplementedCryptoVoteServiceServer) RetrieveAllCryptoVoteByFilter(*RetrieveCryptoReq, CryptoVoteService_RetrieveAllCryptoVoteByFilterServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveAllCryptoVoteByFilter not implemented")
}
func (UnimplementedCryptoVoteServiceServer) UpdateOneCryptoVoteByFilter(context.Context, *UpdateCryptoReq) (*TotalChangesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOneCryptoVoteByFilter not implemented")
}
func (UnimplementedCryptoVoteServiceServer) AddUpVote(context.Context, *AddVoteReq) (*TotalChangesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUpVote not implemented")
}
func (UnimplementedCryptoVoteServiceServer) AddDownVote(context.Context, *AddVoteReq) (*TotalChangesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDownVote not implemented")
}
func (UnimplementedCryptoVoteServiceServer) DeleteAllCryptoVoteByFilter(context.Context, *DeleteCryptoReq) (*TotalChangesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllCryptoVoteByFilter not implemented")
}
func (UnimplementedCryptoVoteServiceServer) mustEmbedUnimplementedCryptoVoteServiceServer() {}

// UnsafeCryptoVoteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CryptoVoteServiceServer will
// result in compilation errors.
type UnsafeCryptoVoteServiceServer interface {
	mustEmbedUnimplementedCryptoVoteServiceServer()
}

func RegisterCryptoVoteServiceServer(s grpc.ServiceRegistrar, srv CryptoVoteServiceServer) {
	s.RegisterService(&CryptoVoteService_ServiceDesc, srv)
}

func _CryptoVoteService_CreateCryptoVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCryptoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).CreateCryptoVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/CreateCryptoVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).CreateCryptoVote(ctx, req.(*CreateCryptoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_RetrieveAllCryptoVoteByFilter_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RetrieveCryptoReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CryptoVoteServiceServer).RetrieveAllCryptoVoteByFilter(m, &cryptoVoteServiceRetrieveAllCryptoVoteByFilterServer{stream})
}

type CryptoVoteService_RetrieveAllCryptoVoteByFilterServer interface {
	Send(*CryptoVote) error
	grpc.ServerStream
}

type cryptoVoteServiceRetrieveAllCryptoVoteByFilterServer struct {
	grpc.ServerStream
}

func (x *cryptoVoteServiceRetrieveAllCryptoVoteByFilterServer) Send(m *CryptoVote) error {
	return x.ServerStream.SendMsg(m)
}

func _CryptoVoteService_UpdateOneCryptoVoteByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCryptoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).UpdateOneCryptoVoteByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/UpdateOneCryptoVoteByFilter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).UpdateOneCryptoVoteByFilter(ctx, req.(*UpdateCryptoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_AddUpVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddVoteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).AddUpVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/AddUpVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).AddUpVote(ctx, req.(*AddVoteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_AddDownVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddVoteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).AddDownVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/AddDownVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).AddDownVote(ctx, req.(*AddVoteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_DeleteAllCryptoVoteByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCryptoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).DeleteAllCryptoVoteByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/DeleteAllCryptoVoteByFilter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).DeleteAllCryptoVoteByFilter(ctx, req.(*DeleteCryptoReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CryptoVoteService_ServiceDesc is the grpc.ServiceDesc for CryptoVoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CryptoVoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cryptovotepb.CryptoVoteService",
	HandlerType: (*CryptoVoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCryptoVote",
			Handler:    _CryptoVoteService_CreateCryptoVote_Handler,
		},
		{
			MethodName: "UpdateOneCryptoVoteByFilter",
			Handler:    _CryptoVoteService_UpdateOneCryptoVoteByFilter_Handler,
		},
		{
			MethodName: "AddUpVote",
			Handler:    _CryptoVoteService_AddUpVote_Handler,
		},
		{
			MethodName: "AddDownVote",
			Handler:    _CryptoVoteService_AddDownVote_Handler,
		},
		{
			MethodName: "DeleteAllCryptoVoteByFilter",
			Handler:    _CryptoVoteService_DeleteAllCryptoVoteByFilter_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RetrieveAllCryptoVoteByFilter",
			Handler:       _CryptoVoteService_RetrieveAllCryptoVoteByFilter_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "cryptovotepb/go_grpc_cryptovote.proto",
}
