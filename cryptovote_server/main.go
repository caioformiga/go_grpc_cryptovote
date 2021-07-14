package main

// No cryptovote_server/main.go temos um gRPC server
// Implements the CryptoVoteService definnnido em cryptovotepb/go_grpc_cryptovote.proto

// source: https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	pb "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of CryptoVote")
)

type cryptoVoteServiceServer struct {
	pb.UnimplementedCryptoVoteServiceServer

	//retornos do tipo streams definidos no cryptoVoteServiceServer do arquivo go_grpc_cryptovote.proto
	savedCryptoCurrencis []*pb.CryptoCurrency

	// protege cryptovotes de serem computados simultaneamente
	// savedCryptoVotes []*pb.CryptoVote
	// mutex sync.Mutex
}

func (s *cryptoVoteServiceServer) ListAllCryptoCurrencies(empty *pb.EmptyReq, stream pb.CryptoVoteService_ListAllCryptoCurrenciesServer) error {
	for _, cryptoCurrencis := range s.savedCryptoCurrencis {
		if err := stream.Send(cryptoCurrencis); err != nil {
			return err
		}
	}
	// se estiver vazio
	return nil
}

// loadCryptoCurrencisFromMongoDB carrega CryptoCurrencis do MongoDB
func (s *cryptoVoteServiceServer) loadCryptoCurrencisFromMongoDB() {
	/*
		data = recuperar todas as cryptoccurencis do banco

		if err := json.Unmarshal(data, &s.savedCryptoCurrencis); err != nil {
			log.Fatalf("Erro ao tentar recuperar  CryptoCurrencis do MongoDB: %v", err)
		}
	*/
}

func newServer() *cryptoVoteServiceServer {
	s := &cryptoVoteServiceServer{}
	s.loadCryptoCurrencisFromMongoDB()
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	cryptovoteGrpcServer := grpc.NewServer(opts...)

	pb.RegisterCryptoVoteServiceServer(cryptovoteGrpcServer, newServer())

	log.Printf("cryptovote service is a web server using gRPC")
	log.Printf("listening at %v", lis.Addr())

	cryptovoteGrpcServer.Serve(lis)
}
