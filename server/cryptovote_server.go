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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/status"

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
	// criar um filtro bson.M com base no model.CryptoVote para recuperar os dados com base neste filtro
	var modelFilter bson.M
	messageCryptoVote := filter.Crypto
	if messageCryptoVote != nil {
		modelFilter = bson.M{
			"name":          messageCryptoVote.Name,
			"symbol":        messageCryptoVote.Symbol,
			"qtd_upvotes":   messageCryptoVote.QtdUpvote,
			"qtd_downvotes": messageCryptoVote.QtdDownvote,
		}
	}

	modelDataCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(modelFilter)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de []model.CryptoVote para json
	jsonData, err := json.Marshal(modelDataCryptoVotes)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	} else {
		// convertendo de json para []pb.CryptoVote
		err = json.Unmarshal(jsonData, &s.cachedCryptoVotes)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal de json para CryptoVotes: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return status.Errorf(grpc_error_code, error_message)
		}
	}

	if len(s.cachedCryptoVotes) > 0 {
		for _, cryptoVotes := range s.cachedCryptoVotes {
			if err := stream.Send(cryptoVotes); err != nil {
				error_message := err.Error()
				grpc_error_code := status.Code(err)
				return status.Errorf(grpc_error_code, error_message)
			}
		}
	}
	return nil
}

func (s *cryptoVoteServiceServer) CreateCryptoVote(ctx context.Context, messageCryptoVote *pb.CryptoVote) (*pb.CryptoVote, error) {
	var modelCryptoVote model.CryptoVote

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVote)
	if err != nil {
		log.Printf("Problemas para fazer Marshal de pb.CryptoVote para json: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
	} else {
		// convertendo de json para model.CryptoVote
		err := json.Unmarshal(jsonData, &modelCryptoVote)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal de json para model.CryptoVote: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
		}
	}

	modelCryptoVote, err = bo.ValidateCryptoVoteData(modelCryptoVote.Name, modelCryptoVote.Symbol, modelCryptoVote.Qtd_Upvote, modelCryptoVote.Qtd_Downvote)
	if err != nil {
		log.Printf("Problemas para validar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := codes.InvalidArgument
		return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
	}

	modelCryptoVote, err = bo.CreateCryptoVote(modelCryptoVote)
	if err != nil {
		log.Printf("Problemas para salvar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de model.CryptoVote para json
	jsonData, err = json.Marshal(modelCryptoVote)
	if err != nil {
		log.Printf("Problemas para fazer Marshal de model.CryptoVote para json: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
	} else {
		// convertendo de json para []pb.CryptoVote
		err := json.Unmarshal(jsonData, &messageCryptoVote)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal de json para []pb.CryptoVote: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return messageCryptoVote, status.Errorf(grpc_error_code, error_message)
		}
	}
	return messageCryptoVote, err
}

func (s *cryptoVoteServiceServer) UpdateAllCryptoVoteByFilter(ctx context.Context, filter *pb.FilterCryptoVote) (*pb.CryptoVote, error) {
	log.Printf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

func (s *cryptoVoteServiceServer) DeleteAllCryptoVoteByFilter(ctx context.Context, filter *pb.FilterCryptoVote) (*pb.CryptoVote, error) {
	log.Printf("Função não implementada ainda...retornando o próprio filtro: %v", filter.Crypto)
	return filter.Crypto, nil
}

// loadCryptoVotesFromMongoDB carrega CryptoVotes
func (s *cryptoVoteServiceServer) loadCryptoVotes() {
	var modelDataCryptoVotes []model.CryptoVote

	// filtro vazio para recuperar todos os dados do banco
	filter := bson.M{}
	modelDataCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(filter)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
	}
	if len(modelDataCryptoVotes) > 0 {

		// convertendo de []model.CryptoVote para json
		jsonData, err := json.Marshal(modelDataCryptoVotes)
		if err != nil {
			log.Printf("Problemas para recuperar dados: %v", err)
		}

		// convertendo de json para []pb.CryptoVote
		err = json.Unmarshal(jsonData, &s.cachedCryptoVotes)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal de json para []pb.CryptoVote: %v", err)
		}

	} else if len(jsonDefalutDataCryptoVotes) > 0 {
		// convertendo de json para []model.CryptoVote
		err := json.Unmarshal(jsonDefalutDataCryptoVotes, &modelDataCryptoVotes)
		if err != nil {
			log.Printf("Erro ao fazer Unmarshal de json para []model.CryptoVote: %v", err)
		}

		// depois do Unmarshal se tiver CryptoVote faz um create
		if len(modelDataCryptoVotes) > 0 {
			for _, modelCryptoVote := range modelDataCryptoVotes {

				validatedCryptoVote, err := bo.ValidateCryptoVoteData(modelCryptoVote.Name, modelCryptoVote.Symbol, modelCryptoVote.Qtd_Upvote, modelCryptoVote.Qtd_Downvote)
				if err != nil {
					log.Printf("Problemas para validar dados: %v", err)
				}

				modelCryptoVote, err = bo.CreateCryptoVote(validatedCryptoVote)
				if err != nil {
					log.Printf("Problemas para salvar dados: %v", err)
				}
			}
		}
		//chama recursivamente para entrar no primeiro laço
		s.loadCryptoVotes()
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

	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
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
			log.Printf("Failed to generate credentials %v", err)
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
