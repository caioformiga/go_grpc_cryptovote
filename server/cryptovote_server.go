package main

// No cryptovote_server/main.go temos um gRPC server
// Implements the CryptoVoteService definnnido em cryptovotepb/go_grpc_cryptovote.proto

// source: https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/reflection"

	pb "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb"
)

/*
TODO usar env pois quaando for rodar em kubernetes fafilita a customização pelo config (yaml)
*/
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
	for _, cryptoCurrency := range s.savedCryptoCurrencis {
		if err := stream.Send(cryptoCurrency); err != nil {
			return err
		}
	}
	// se estiver vazio
	return nil
}

// loadCryptoCurrencisFromMongoDB carrega CryptoCurrencis do MongoDB
func (s *cryptoVoteServiceServer) loadCryptoCurrencis() {
	var data []byte = exampleData
	if err := json.Unmarshal(data, &s.savedCryptoCurrencis); err != nil {
		log.Fatalf("Erro ao tentar recuperar CryptoCurrencis de exampleData: %v", err)
	}
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printAllCrypto() {
	log.Printf("Available Crypto's:")
	for i, cryptoCurrency := range s.savedCryptoCurrencis {
		log.Printf("Crypto index: %d", i)
		log.Printf("Crypto Name: %v", cryptoCurrency.Name)
		log.Printf("Crypto Symbol: %v", cryptoCurrency.Symbol)
		log.Printf("\n")
	}
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printServerInfo() {
	log.Printf("Cryptovote service is a web server using gRPC")
}

func newServer() *cryptoVoteServiceServer {
	s := &cryptoVoteServiceServer{}
	s.printServerInfo()
	s.loadCryptoCurrencis()
	s.printAllCrypto()
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

	// registranndo o servidor como reflection service para usar o gRPCui
	reflection.Register(cryptovoteGrpcServer)

	log.Printf("listening at %v", lis.Addr())

	cryptovoteGrpcServer.Serve(lis)
}

// exampleData is a json object baseado em CryptoCurrency message
/*
   string name = 1;
   string symbol = 2;
*/
var exampleData = []byte(`[{
    "name": "Bitcoin",
    "symbol": "BTC"
}, {
	"name": "Ethereum",
    "symbol": "ETH"
}, {	
	"name": "Klever",
    "symbol": "KLV"
}]`)
