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

	"go.mongodb.org/mongo-driver/bson"

	"github.com/caioformiga/go_mongodb_crud_cryptovote/bo"
	"github.com/caioformiga/go_mongodb_crud_cryptovote/model"

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
	cachedCryptoVotes []*pb.CryptoVote

	// protege cryptovotes de serem computados simultaneamente
	mutex sync.Mutex
}

// Handler dos services definidos em go_grpc_cryptovote_grpc.pb.go
/*
	RetrieveAllCryptoVoteByFilter recupera todos os CryptoVotes encapulados em uma struct pb.CryptoVote
	primeiro tenta recuperar da memoria (cachedCryptoVotes), caso essteja vazio faz uma busca no banco de
	dados e carrega os dado na memoria
*/
func (s *cryptoVoteServiceServer) RetrieveAllCryptoVoteByFilter(filter *pb.FilterCryptoVote, stream pb.CryptoVoteService_RetrieveAllCryptoVoteByFilterServer) error {
	if len(s.cachedCryptoVotes) <= 0 {
		s.loadCryptoVotes()
	}

	for _, cryptoVotes := range s.cachedCryptoVotes {
		if err := stream.Send(cryptoVotes); err != nil {
			return err
		}
	}
	return nil
}

func (s *cryptoVoteServiceServer) CreateCryptoVote(ctx context.Context, messageCryptoVote *pb.CryptoVote) (*pb.CryptoVote, error) {

	modelCryptoVote, err := TranslateToModelCryptoVote(messageCryptoVote)
	if err != nil {
		log.Fatalf("Problemas para traduzir dados: %v", err)
	}

	modelCryptoVote, err = bo.CreateCryptoVote(modelCryptoVote.Name, modelCryptoVote.Symbol, modelCryptoVote.Qtd_Upvote, modelCryptoVote.Qtd_Downvote)
	if err != nil {
		log.Fatalf("Problemas para salvar dados: %v", err)
	}
	log.Printf("modelCryptoVote: %v", modelCryptoVote)

	messageCryptoVote, err = TranslateToMessageCryptoVote(modelCryptoVote)
	if err != nil {
		log.Fatalf("Problemas para traduzir dados: %v", err)
	}
	return messageCryptoVote, err
}

func (s *cryptoVoteServiceServer) UpdateAllCryptoVoteByFilter(ctx context.Context, filter *pb.FilterCryptoVote) (*pb.CryptoVote, error) {
	log.Fatalf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

func (s *cryptoVoteServiceServer) DeleteAllCryptoVoteByFilter(ctx context.Context, filter *pb.FilterCryptoVote) (*pb.CryptoVote, error) {
	log.Fatalf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

// loadCryptoVotesFromMongoDB carrega CryptoVotes
func (s *cryptoVoteServiceServer) loadCryptoVotes() {
	var modelDataCryptoVotes []model.CryptoVote

	// filtro vazio para recuperar todos os dados do bancp
	filter := bson.M{}
	modelDataCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(filter)
	if err != nil {
		log.Fatalf("Problemas para recuperar dados: %v", err)
	}
	if len(modelDataCryptoVotes) > 0 {
		log.Print("Crypto do banco de dados...")

		// convertendo de model.CryptoVote para json
		jsonData, err := json.Marshal(modelDataCryptoVotes)
		if err != nil {
			log.Fatalf("Problemas para fazer Marshal: %v", err)
		} else {
			// convertendo de json para []pb.CryptoVote
			err := json.Unmarshal(jsonData, &s.cachedCryptoVotes)
			if err != nil {
				log.Fatalf("Erro ao fazer Unmarshal de CryptoVotes em cachedCryptoVotes: %v", err)
			}
		}
	} else {
		log.Print("Crypto do arquivo jsonDefalutDataCryptoVotes...")

		if len(jsonDefalutDataCryptoVotes) > 0 {
			// convertendo de json para []model.CryptoVote
			err := json.Unmarshal(jsonDefalutDataCryptoVotes, &modelDataCryptoVotes)
			if err != nil {
				log.Fatalf("Erro ao fazer Unmarshal de CryptoVotes em modelDataCryptoVotes: %v", err)
			} else {
				// depois do Unmarshal se tiver CryptoVote
				// para cada CryptoVote faz um create
				if len(modelDataCryptoVotes) > 0 {
					for _, modelCryptoVote := range modelDataCryptoVotes {

						name := modelCryptoVote.Name
						symbol := modelCryptoVote.Symbol
						qtd_upvote := modelCryptoVote.Qtd_Upvote
						qtd_downvote := modelCryptoVote.Qtd_Downvote

						modelCryptoVote, err = bo.CreateCryptoVote(name, symbol, qtd_upvote, qtd_downvote)
						if err != nil {
							log.Fatalf("Problemas para salvar dados: %v", err)
						}
					}
				}
			}
		}
		// convertendo de json para []pb.CryptoVote
		if err := json.Unmarshal(jsonDefalutDataCryptoVotes, &s.cachedCryptoVotes); err != nil {
			log.Fatalf("Erro ao fazer Unmarshal de CryptoVotes em cachedCryptoVotes: %v", err)
		}
	}
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printAllCryptoVotes() {
	log.Printf("Available Crypto's:")
	for i, cryptoVote := range s.cachedCryptoVotes {
		log.Printf("Crypto index: %d", i)
		log.Printf("Crypto Name: %v", cryptoVote.Name)
		log.Printf("Crypto Symbol: %v", cryptoVote.Symbol)
		log.Printf("Crypto QTD Up: %v", cryptoVote.QtdUpvote)
		log.Printf("Crypto QTD Dow: %v", cryptoVote.QtdDownvote)
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
	s.loadCryptoVotes()
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

// jsonDefalutDataCryptoVotes é um objeto json baseado em CryptoVote message
/*
    string name = 1;
	string symbol = 2;
	int32 qtd_upvote = 3;
	int32 qtd_downvote = 4;
*/
var jsonDefalutDataCryptoVotes = []byte(`[{
	"name": "Bitcoin",
    "symbol": "BTC",
	"qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"name": "Ethereum",
    "symbol": "ETH",
    "qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"name": "Klever",
	"symbol": "KLV",
    "qtd_upvote": 0,
	"qtd_downvote": 0	
}]`)
