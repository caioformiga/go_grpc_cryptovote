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
	// Obtains all CryptoCurrency
	ListAllCryptoCurrencies(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetCryptocurrencyResponse, error)
	GetStreamCryptoCurrencies(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (CryptoVoteService_GetStreamCryptoCurrenciesClient, error)
	// Obtains all CryptoVote
	ListAllCryptoVotes(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetCryptoVoteResponse, error)
	GetStreamCryptoVotes(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (CryptoVoteService_GetStreamCryptoVotesClient, error)
}

type cryptoVoteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCryptoVoteServiceClient(cc grpc.ClientConnInterface) CryptoVoteServiceClient {
	return &cryptoVoteServiceClient{cc}
}

func (c *cryptoVoteServiceClient) ListAllCryptoCurrencies(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetCryptocurrencyResponse, error) {
	out := new(GetCryptocurrencyResponse)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/ListAllCryptoCurrencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) GetStreamCryptoCurrencies(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (CryptoVoteService_GetStreamCryptoCurrenciesClient, error) {
	stream, err := c.cc.NewStream(ctx, &CryptoVoteService_ServiceDesc.Streams[0], "/cryptovotepb.CryptoVoteService/GetStreamCryptoCurrencies", opts...)
	if err != nil {
		return nil, err
	}
	x := &cryptoVoteServiceGetStreamCryptoCurrenciesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CryptoVoteService_GetStreamCryptoCurrenciesClient interface {
	Recv() (*Cryptocurrency, error)
	grpc.ClientStream
}

type cryptoVoteServiceGetStreamCryptoCurrenciesClient struct {
	grpc.ClientStream
}

func (x *cryptoVoteServiceGetStreamCryptoCurrenciesClient) Recv() (*Cryptocurrency, error) {
	m := new(Cryptocurrency)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cryptoVoteServiceClient) ListAllCryptoVotes(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (*GetCryptoVoteResponse, error) {
	out := new(GetCryptoVoteResponse)
	err := c.cc.Invoke(ctx, "/cryptovotepb.CryptoVoteService/ListAllCryptoVotes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cryptoVoteServiceClient) GetStreamCryptoVotes(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (CryptoVoteService_GetStreamCryptoVotesClient, error) {
	stream, err := c.cc.NewStream(ctx, &CryptoVoteService_ServiceDesc.Streams[1], "/cryptovotepb.CryptoVoteService/GetStreamCryptoVotes", opts...)
	if err != nil {
		return nil, err
	}
	x := &cryptoVoteServiceGetStreamCryptoVotesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CryptoVoteService_GetStreamCryptoVotesClient interface {
	Recv() (*CryptoVote, error)
	grpc.ClientStream
}

type cryptoVoteServiceGetStreamCryptoVotesClient struct {
	grpc.ClientStream
}

func (x *cryptoVoteServiceGetStreamCryptoVotesClient) Recv() (*CryptoVote, error) {
	m := new(CryptoVote)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CryptoVoteServiceServer is the server API for CryptoVoteService service.
// All implementations must embed UnimplementedCryptoVoteServiceServer
// for forward compatibility
type CryptoVoteServiceServer interface {
	// Obtains all CryptoCurrency
	ListAllCryptoCurrencies(context.Context, *EmptyReq) (*GetCryptocurrencyResponse, error)
	GetStreamCryptoCurrencies(*EmptyReq, CryptoVoteService_GetStreamCryptoCurrenciesServer) error
	// Obtains all CryptoVote
	ListAllCryptoVotes(context.Context, *EmptyReq) (*GetCryptoVoteResponse, error)
	GetStreamCryptoVotes(*EmptyReq, CryptoVoteService_GetStreamCryptoVotesServer) error
	mustEmbedUnimplementedCryptoVoteServiceServer()
}

// UnimplementedCryptoVoteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCryptoVoteServiceServer struct {
}

func (UnimplementedCryptoVoteServiceServer) ListAllCryptoCurrencies(context.Context, *EmptyReq) (*GetCryptocurrencyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllCryptoCurrencies not implemented")
}
func (UnimplementedCryptoVoteServiceServer) GetStreamCryptoCurrencies(*EmptyReq, CryptoVoteService_GetStreamCryptoCurrenciesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStreamCryptoCurrencies not implemented")
}
func (UnimplementedCryptoVoteServiceServer) ListAllCryptoVotes(context.Context, *EmptyReq) (*GetCryptoVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAllCryptoVotes not implemented")
}
func (UnimplementedCryptoVoteServiceServer) GetStreamCryptoVotes(*EmptyReq, CryptoVoteService_GetStreamCryptoVotesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStreamCryptoVotes not implemented")
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

func _CryptoVoteService_ListAllCryptoCurrencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).ListAllCryptoCurrencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/ListAllCryptoCurrencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).ListAllCryptoCurrencies(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_GetStreamCryptoCurrencies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CryptoVoteServiceServer).GetStreamCryptoCurrencies(m, &cryptoVoteServiceGetStreamCryptoCurrenciesServer{stream})
}

type CryptoVoteService_GetStreamCryptoCurrenciesServer interface {
	Send(*Cryptocurrency) error
	grpc.ServerStream
}

type cryptoVoteServiceGetStreamCryptoCurrenciesServer struct {
	grpc.ServerStream
}

func (x *cryptoVoteServiceGetStreamCryptoCurrenciesServer) Send(m *Cryptocurrency) error {
	return x.ServerStream.SendMsg(m)
}

func _CryptoVoteService_ListAllCryptoVotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CryptoVoteServiceServer).ListAllCryptoVotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cryptovotepb.CryptoVoteService/ListAllCryptoVotes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CryptoVoteServiceServer).ListAllCryptoVotes(ctx, req.(*EmptyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CryptoVoteService_GetStreamCryptoVotes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CryptoVoteServiceServer).GetStreamCryptoVotes(m, &cryptoVoteServiceGetStreamCryptoVotesServer{stream})
}

type CryptoVoteService_GetStreamCryptoVotesServer interface {
	Send(*CryptoVote) error
	grpc.ServerStream
}

type cryptoVoteServiceGetStreamCryptoVotesServer struct {
	grpc.ServerStream
}

func (x *cryptoVoteServiceGetStreamCryptoVotesServer) Send(m *CryptoVote) error {
	return x.ServerStream.SendMsg(m)
}

// CryptoVoteService_ServiceDesc is the grpc.ServiceDesc for CryptoVoteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CryptoVoteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cryptovotepb.CryptoVoteService",
	HandlerType: (*CryptoVoteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAllCryptoCurrencies",
			Handler:    _CryptoVoteService_ListAllCryptoCurrencies_Handler,
		},
		{
			MethodName: "ListAllCryptoVotes",
			Handler:    _CryptoVoteService_ListAllCryptoVotes_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetStreamCryptoCurrencies",
			Handler:       _CryptoVoteService_GetStreamCryptoCurrencies_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStreamCryptoVotes",
			Handler:       _CryptoVoteService_GetStreamCryptoVotes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "cryptovotepb/go_grpc_cryptovote.proto",
}
