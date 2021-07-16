package main

// No cryptovote_server/main.go temos um gRPC server
// Implements the CryptoVoteService definnnido em cryptovotepb/go_grpc_cryptovote.proto

// source: https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/reflection"

	bo "github.com/caioformiga/go_mongodb_crud_cryptovote/bo"

	pb "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb"
)

/*
TODO usar env pois quando for rodar em kubernetes fafilita a customização pelo config (yaml)
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

	//retornos do tipo streams definidos no cryptoVoteServiceServer do arquivo go_grpc_cryptovote.proto
	savedCryptoVotes []*pb.CryptoVote

	// protege cryptovotes de serem computados simultaneamente
	mutex sync.Mutex
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

func (s *cryptoVoteServiceServer) ListAllCryptoVotes(empty *pb.EmptyReq, stream pb.CryptoVoteService_ListAllCryptoVotesServer) error {
	for _, cryptoVotes := range s.savedCryptoVotes {
		if err := stream.Send(cryptoVotes); err != nil {
			return err
		}
	}
	// se estiver vazio
	return nil
}

func (s *cryptoVoteServiceServer) CreateCryptoCurrency(ctx context.Context, newCryptoCurrency *pb.CryptoCurrency) (*pb.CryptoCurrency, error) {
	mongoCryptoCurrency, err := bo.CreateCryptoCurrency(newCryptoCurrency.Name, newCryptoCurrency.Symbol)
	if err != nil {
		log.Fatalf("Problemas para salvar dados: %v", err)
	}
	log.Printf("mongoCryptoCurrency: %v", mongoCryptoCurrency)
	return newCryptoCurrency, err
}

func (s *cryptoVoteServiceServer) RetrieveAllCryptoCurrencyByFilter(ctx context.Context, filter *pb.FilterCryptoCurrency) (*pb.CryptoCurrency, error) {
	log.Fatalf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

func (s *cryptoVoteServiceServer) UpdateAllCryptoCurrencyByFilter(ctx context.Context, filter *pb.FilterCryptoCurrency) (*pb.CryptoCurrency, error) {
	log.Fatalf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

func (s *cryptoVoteServiceServer) DeleteAllCryptoCurrencyByFilter(ctx context.Context, filter *pb.FilterCryptoCurrency) (*pb.CryptoCurrency, error) {
	log.Fatalf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

// loadCryptoCurrencisFromMongoDB carrega CryptoCurrencis
func (s *cryptoVoteServiceServer) loadCryptoCurrencis() {
	var data []byte = exampleDataCryptoCurrencies
	if err := json.Unmarshal(data, &s.savedCryptoCurrencis); err != nil {
		log.Fatalf("Erro ao tentar recuperar CryptoCurrencis de exampleDataCryptoCurrencies: %v", err)
	}
}

// loadCryptoCurrencisFromMongoDB carrega CryptoVotes
func (s *cryptoVoteServiceServer) loadCryptoVotes() {
	var data []byte = exampleDataCryptoVotes
	if err := json.Unmarshal(data, &s.savedCryptoVotes); err != nil {
		log.Fatalf("Erro ao tentar recuperar CryptoCurrencis de exampleDataCryptoVotes: %v", err)
	}
}

// imprime no console todos os campos do
func printCryptoVotesToString(cryptoVote *pb.CryptoVote) {
	log.Printf("Crypto Name: %v", cryptoVote.Crypto.Name)
	log.Printf("Crypto Symbol: %v", cryptoVote.Crypto.Symbol)
	log.Printf("Crypto QTD Up: %v", cryptoVote.QtdUpvote)
	log.Printf("Crypto QTD Dow: %v", cryptoVote.QtdDownvote)
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printAllCryptoVotes() {
	log.Printf("Available Crypto's:")
	for i, cryptoVote := range s.savedCryptoVotes {
		log.Printf("Crypto index: %d", i)
		printCryptoVotesToString(cryptoVote)
		log.Printf("\n")
	}
}

// imprime no console todos os campos do
func printCryptoCurrencisToString(cryptoCurrency *pb.CryptoCurrency) {
	log.Printf("Crypto Name: %v", cryptoCurrency.Name)
	log.Printf("Crypto Symbol: %v", cryptoCurrency.Symbol)
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printAllCryptoCurrencis() {
	log.Printf("Available Crypto's:")
	for i, cryptoCurrency := range s.savedCryptoCurrencis {
		log.Printf("Crypto index: %d", i)
		printCryptoCurrencisToString(cryptoCurrency)
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
	s.loadCryptoVotes()
	s.printAllCryptoCurrencis()
	s.printAllCryptoVotes()
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

	// registranndo o servidor como reflection service para usar o gRPCuic
	reflection.Register(cryptovoteGrpcServer)

	log.Printf("listening at %v", lis.Addr())

	cryptovoteGrpcServer.Serve(lis)
}

// exampleDataCryptoCurrencies é um objeto json baseado em CryptoCurrency message
/*
   string name = 1;
   string symbol = 2;
*/
var exampleDataCryptoCurrencies = []byte(`[{
    "name": "Bitcoin",
    "symbol": "BTC"
}, {
	"name": "Ethereum",
    "symbol": "ETH"
}, {	
	"name": "Klever",
    "symbol": "KLV"
}]`)

// exampleDataCryptoVote é um objeto json baseado em CryptoVotes message
/*
   CryptoCurrency crypto = 1;
   int32 qtd_upvote = 2;
   int32 qtd_downvote = 3;
   string idHex = 4;
*/
var exampleDataCryptoVotes = []byte(`[{
	"crypto": {
		"name": "Bitcoin",
    	"symbol": "BTC"
	},
    "qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"crypto": {
		"name": "Ethereum",
    	"symbol": "ETH"
	},
    "qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"crypto": {
		"name": "Klever",
		"symbol": "KLV"
	},
    "qtd_upvote": 0,
	"qtd_downvote": 0	
}]`)
