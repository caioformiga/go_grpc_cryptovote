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
	CreateCryptoVote faz Marshal/Unmarshal do buffer para JSON para CryptoVote e faz validação antes de
	salvar no banco
*/
func (s *cryptoVoteServiceServer) CreateCryptoVote(ctx context.Context, req *pb.CreateCryptoReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Crypto

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// cria a struct que armazena os dados do json
	var cyptoVote model.CryptoVote

	// convertendo de json para model.CryptoVote
	err = json.Unmarshal(jsonData, &cyptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// faz as validações e salva no banco
	// validação de valores default
	// validação de unicos para name e symbol
	cyptoVote, err = bo.CreateCryptoVote(cyptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// registra o total de atualizações
	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

/*
	RetrieveAllCryptoVoteByFilter recupera todas as CryptoVotes com base no filtro em uma struct pb.CryptoVote
	primeiro tenta recuperar da memoria (cachedCryptoVotes), caso essteja vazio faz uma busca no banco de
	dados e carrega os dado na memoria
*/
func (s *cryptoVoteServiceServer) RetrieveAllCryptoVoteByFilter(req *pb.RetrieveCryptoReq, stream pb.CryptoVoteService_RetrieveAllCryptoVoteByFilterServer) error {
	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVote)
	if err != nil {
		log.Printf("Problemas para fazer Marshal de pb.CryptoVote para json: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de json para bson
	var filter = bson.M{}
	err = json.Unmarshal(jsonData, &filter)
	if err != nil {
		fmt.Printf("Problemas para fazer Unmarshal de json para bson.M: %v", err)
	}

	modelCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(filter)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de []model.CryptoVote para json
	jsonData, err = json.Marshal(modelCryptoVotes)
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

/*
	UpdateOneCryptoVoteByFilter atualiza um CryptoVote campos nome e symbol om base no filtro
*/
func (s *cryptoVoteServiceServer) UpdateOneCryptoVoteByFilter(ctx context.Context, req *pb.UpdateCryptoReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVoteFilter := req.Filter

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVoteFilter)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de json para bson
	var filter = bson.M{}
	err = json.Unmarshal(jsonData, &filter)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// cria um model.CryptoVote com dados da req
	var cryptoVoteNewData = model.CryptoVote{
		Name:   req.NewName,
		Symbol: req.NewSymbol,
	}

	// se exister um único registro com base no filtro faz uma atualização
	// no name e symbol de uma CryptoVote, antes faz validação e salva
	// validação de valores default
	// validação de unicos para name e symbol
	cyptoVote, err := bo.UpdateOneCryptoVoteByFilter(filter, cryptoVoteNewData)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

func (s *cryptoVoteServiceServer) AddUpVote(ctx context.Context, req *pb.AddVoteReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	// cria um model.CryptoVote com dados da req
	cyptoVote, err := bo.AddUpVote(req.Name, req.Symbol)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// se o ID não for zero sinal que conseguiu realizar UpVote
	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

func (s *cryptoVoteServiceServer) AddDownVote(ctx context.Context, req *pb.AddVoteReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	// cria um model.CryptoVote com dados da req
	cyptoVote, err := bo.AddDownVote(req.Name, req.Symbol)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// se o ID não for zero sinal que conseguiu realizar DownVote
	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

func (s *cryptoVoteServiceServer) DeleteAllCryptoVoteByFilter(ctx context.Context, req *pb.DeleteCryptoReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	// retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de json para bson
	var filter = bson.M{}
	err = json.Unmarshal(jsonData, &filter)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	deletedQtd, err := bo.DeleteAllCryptoVoteByFilter(filter)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	if deletedQtd > 0 {
		res.Qtd = int32(deletedQtd)
	}
	return res, err
}

// loadCryptoVotesFromMongoDB carrega CryptoVotes
func (s *cryptoVoteServiceServer) loadCryptoVotes() {
	var modelCryptoVotes []model.CryptoVote

	// filtro vazio para recuperar todos os dados do banco
	filter := bson.M{}
	modelCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(filter)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		err = status.Errorf(grpc_error_code, error_message)
		log.Printf("[server.cryptovote_server.go] +%v", err)
	}
	if len(modelCryptoVotes) > 0 {

		// convertendo de []model.CryptoVote para json
		jsonData, err := json.Marshal(modelCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[server.cryptovote_server.go] +%v", err)
		}

		// convertendo de json para []pb.CryptoVote
		err = json.Unmarshal(jsonData, &s.cachedCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[server.cryptovote_server.go] +%v", err)
		}

	} else if len(jsonDefalutDataCryptoVotes) > 0 {
		// convertendo de json para []model.CryptoVote
		err := json.Unmarshal(jsonDefalutDataCryptoVotes, &modelCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[server.cryptovote_server.go] +%v", err)
		}

		// depois do Unmarshal se tiver CryptoVote faz um create
		if len(modelCryptoVotes) > 0 {
			for _, modelCryptoVote := range modelCryptoVotes {

				modelCryptoVote, err = bo.CreateCryptoVote(modelCryptoVote)
				if err != nil {
					error_message := err.Error()
					grpc_error_code := status.Code(err)
					err = status.Errorf(grpc_error_code, error_message)
					log.Printf("[server.cryptovote_server.go] +%v", err)
				}
			}
		}
		//chama recursivamente para entrar no primeiro laço
		s.loadCryptoVotes()
	}
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printServerInfo() {
	log.Printf("[server.cryptovote_server.go] Cryptovote service is a web server using gRPC")
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
		log.Printf("[server.cryptovote_server.go] failed to listen: %v", err)
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
			log.Printf("[server.cryptovote_server.go] Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	cryptovoteGrpcServer := grpc.NewServer(opts...)
	pb.RegisterCryptoVoteServiceServer(cryptovoteGrpcServer, newServer())

	// registranndo o servidor como reflection service para usar o gRPCuic
	reflection.Register(cryptovoteGrpcServer)

	log.Printf("[server.cryptovote_server.go] listening at %v", lis.Addr())

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
