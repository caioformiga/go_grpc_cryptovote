package main

// No cryptovote_server/main.go temos um gRPC server
// Implements the CryptoVoteService definnnido em cryptovotepb/go_grpc_cryptovote.proto

// source: https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/status"

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
	port     = flag.Int("port", 50000, "The server port")
	//jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of CryptoVote")
)

type cryptoVoteServiceServer struct {
	pb.UnimplementedCryptoVoteServiceServer

	//retornos do tipo streams definidos no cryptoVoteServiceServer do arquivo go_grpc_cryptovote.proto
	cachedCryptoVotes []*pb.CryptoVote

	// protege alterações de escrita serem realizadas simultaneamente
	mutex sync.Mutex
}

func validateArgumentsType(messageCryptoVote *pb.CryptoVote) (bool, error) {
	var validate bool = false

	// se messageCryptoVote não for diferente de nulo
	// sinal que todos os argumentos serão nulos
	// retorna true pois não tem o que validar
	if !(messageCryptoVote != nil) {
		return true, nil
	}

	// quando usamos reflect.TypeOf(interface) a fun String() retorna o tipo em formato string
	// esse mapeamento esta na var kindNames de type.go
	if !(reflect.TypeOf(messageCryptoVote.Name).String() == "string") {
		z := "[cryptovote.grpc_server] dado informado para symbol (" + messageCryptoVote.Symbol + ") não pode ser diferente de string"
		return false, errors.New(z)
	} else {
		validate = true
	}

	if !(reflect.TypeOf(messageCryptoVote.Symbol).String() == "string") {
		z := "[cryptovote.grpc_server] dado informado para symbol (" + messageCryptoVote.Symbol + ") não pode ser diferente de string"
		return false, errors.New(z)
	} else {
		validate = true
	}

	if !(reflect.TypeOf(messageCryptoVote.QtdUpvote).String() == "int32") {
		z := "[cryptovote.grpc_server] dado informado para qtd_upvote (" + string(messageCryptoVote.QtdUpvote) + ") não pode ser diferente de int32"
		return false, errors.New(z)
	} else {
		validate = true
	}

	if !(reflect.TypeOf(messageCryptoVote.QtdDownvote).String() == "int32") {
		z := "[cryptovote.grpc_server] dado informado para qtd_downvote (" + string(messageCryptoVote.QtdDownvote) + ") não pode ser diferente de int32"
		return false, errors.New(z)
	} else {
		validate = true
	}

	return validate, nil
}

// Handler do service definidos em go_grpc_cryptovote_grpc.pb.go
/*
	CreateCryptoVote faz Marshal/Unmarshal do buffer para JSON para CryptoVote e faz validação antes de
	salvar no banco

	url de acesso definida via google.api.http
	post: "/v1/cryptovotes"
    body: "crypto"
*/
func (s *cryptoVoteServiceServer) InsertCryptoVote(ctx context.Context, req *pb.InsertCryptoReq) (*pb.TotalChangesRes, error) {
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Crypto

	//populate crypo
	messageCryptoVote.Name = strings.Title(strings.ToLower(strings.TrimSpace(messageCryptoVote.Name)))
	messageCryptoVote.Symbol = strings.ToUpper(strings.TrimSpace(messageCryptoVote.Symbol))

	// faz validação das entradas
	validate, err := validateArgumentsType(messageCryptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	if validate {
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

		s.mutex.Lock()
		// faz as validações e salva no banco
		// validação de valores default
		// validação de unicos para name e symbol
		cyptoVote, err = bo.CreateCryptoVote(cyptoVote)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return res, status.Errorf(grpc_error_code, error_message)
		}
		s.mutex.Unlock()

		// registra o total de atualizações
		if !cyptoVote.Id.IsZero() {
			res.Qtd = 1
		}
	}
	return res, err
}

func helpperPrepareFilter(messageCryptoVote *pb.CryptoVote) (model.CryptoVote, error) {
	var filterCryptoVote model.CryptoVote
	var validate bool = false
	var err error
	var jsonData []byte

	//populate crypo
	messageCryptoVote.Name = strings.Title(strings.ToLower(strings.TrimSpace(messageCryptoVote.Name)))
	messageCryptoVote.Symbol = strings.ToUpper(strings.TrimSpace(messageCryptoVote.Symbol))

	// faz validação das entradas, se ms for nil retorna true pois não tem argumentos para validar
	validate, err = validateArgumentsType(messageCryptoVote)
	if err != nil {
		return filterCryptoVote, err
	}

	if validate {
		// convertendo de pb.CryptoVote para json
		jsonData, err = json.Marshal(messageCryptoVote)
		if err != nil {
			return filterCryptoVote, err
		}

		// convertendo de json para model.CryptoVote
		err = json.Unmarshal(jsonData, &filterCryptoVote)
		if err != nil {
			return filterCryptoVote, err
		}
	}
	return filterCryptoVote, nil
}

/*
	ListAllCryptoVoteByFilter recupera todas as CryptoVotes com base no filtro em uma struct pb.CryptoVote
	primeiro tenta recuperar da memoria (cachedCryptoVotes), caso essteja vazio faz uma busca no banco de
	dados e carrega os dado na memoria

	url de acesso definida via google.api.http
    get: "/v1/cryptovotes/{filter}"
*/
func (s *cryptoVoteServiceServer) ListAllCryptoVoteByFilter(req *pb.ListCryptoReq, stream pb.CryptoVoteService_ListAllCryptoVoteByFilterServer) error {
	var retrievedCryptoVotes []model.CryptoVote
	var filterCryptoVote model.CryptoVote
	var err error
	var jsonData []byte

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter
	if messageCryptoVote != nil {
		filterCryptoVote, err = helpperPrepareFilter(messageCryptoVote)
		if err != nil {
			log.Printf("Problemas para recuperar dados: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return status.Errorf(grpc_error_code, error_message)
		}
	} else {
		// se for nil precisamos criar um filtro vazio
		filterCryptoVote = model.CryptoVote{
			Name:   "",
			Symbol: "",
		}
	}

	retrievedCryptoVotes, err = bo.RetrieveAllCryptoVoteByFilter(filterCryptoVote)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de []model.CryptoVote para json
	jsonData, err = json.Marshal(retrievedCryptoVotes)
	if err != nil {
		log.Printf("Problemas para recuperar dados: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
	}

	// convertendo de json para []pb.CryptoVote
	err = json.Unmarshal(jsonData, &s.cachedCryptoVotes)
	if err != nil {
		log.Printf("Erro ao fazer Unmarshal de json para CryptoVotes: %v", err)
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return status.Errorf(grpc_error_code, error_message)
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
	UpdateOneCryptoVoteByFilter atualiza um CryptoVote usando um filtro

	url de acesso definida via google.api.http
    post: "/v1/cryptovotes"
    body: "*"

*/
func (s *cryptoVoteServiceServer) UpdateOneCryptoVoteByFilter(ctx context.Context, req *pb.UpdateCryptoReq) (*pb.TotalChangesRes, error) {
	var filterCryptoVote model.CryptoVote
	var err error
	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter
	if messageCryptoVote != nil {
		filterCryptoVote, err = helpperPrepareFilter(messageCryptoVote)
		if err != nil {
			log.Printf("Problemas com helpperPrepareFilter: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return res, status.Errorf(grpc_error_code, error_message)
		}
	} else {
		// se for nil precisamos retornar pois não pode loclaizar registro
		log.Printf("Problemas com o filtro de busca: %v", err)
		err = errors.New("[cryptovote.grpc_server] filtro não pode ser empty")
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	// cria um model.CryptoVote com dados da req
	var cryptoVoteNewData = model.CryptoVote{
		Name:   req.NewName,
		Symbol: req.NewSymbol,
	}

	//populate crypo
	cryptoVoteNewData.Name = strings.Title(strings.ToLower(strings.TrimSpace(cryptoVoteNewData.Name)))
	cryptoVoteNewData.Symbol = strings.ToUpper(strings.TrimSpace(cryptoVoteNewData.Symbol))

	s.mutex.Lock()
	// se exister um único registro com base no filtro faz uma atualização
	// no name e symbol de uma CryptoVote, antes faz validação e salva
	// validação de valores default
	// validação de unicos para name e symbol
	cyptoVote, err := bo.UpdateOneCryptoVoteByFilter(filterCryptoVote, cryptoVoteNewData)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}
	s.mutex.Unlock()

	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}

	return res, err
}

/*
	DeleteAllCryptoVoteByFilter remove um CryptoVote usando um filtro

	url de acesso definida via google.api.http
    delete: "/v1/cryptovotes/{filter}"

*/
func (s *cryptoVoteServiceServer) DeleteAllCryptoVoteByFilter(ctx context.Context, req *pb.DeleteCryptoReq) (*pb.TotalChangesRes, error) {
	var filterCryptoVote model.CryptoVote
	var err error

	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter
	if messageCryptoVote != nil {
		filterCryptoVote, err = helpperPrepareFilter(messageCryptoVote)
		if err != nil {
			log.Printf("Problemas com helpperPrepareFilter: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return res, status.Errorf(grpc_error_code, error_message)
		}
	} else {
		// se for nil precisamos retornar pois não pode loclaizar registro
		log.Printf("Problemas com o filtro de busca: %v", err)
		err = errors.New("[cryptovote.grpc_server] filtro não pode ser empty")
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	s.mutex.Lock()
	deletedQtd, err := bo.DeleteAllCryptoVoteByFilter(filterCryptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}
	s.mutex.Unlock()

	if deletedQtd > 0 {
		res.Qtd = int32(deletedQtd)
	}
	return res, err
}

/*
	AddUpVote adiciona um voto em UpVote da CryptoVote do filtro

	url de acesso definida via google.api.http
    post: "/v1:upvote"
     ody: "*"

*/
func (s *cryptoVoteServiceServer) AddUpVote(ctx context.Context, req *pb.AddVoteReq) (*pb.TotalChangesRes, error) {
	var filterCryptoVote model.CryptoVote
	var err error

	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter
	if messageCryptoVote != nil {
		filterCryptoVote, err = helpperPrepareFilter(messageCryptoVote)
		if err != nil {
			log.Printf("Problemas com helpperPrepareFilter: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return res, status.Errorf(grpc_error_code, error_message)
		}
	} else {
		// se for nil precisamos retornar pois não pode loclaizar registro
		log.Printf("Problemas com o filtro de busca: %v", err)
		err = errors.New("[cryptovote.grpc_server] filtro não pode ser empty")
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	s.mutex.Lock()
	// cria um model.CryptoVote com dados da req
	cyptoVote, err := bo.AddUpVote(filterCryptoVote.Name, filterCryptoVote.Symbol)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}
	s.mutex.Unlock()

	// se o ID não for zero sinal que conseguiu realizar UpVote
	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

/*
	AddDownVote adiciona um voto em DownVote da CryptoVote do filtro

	url de acesso definida via google.api.http
    post: "/v1:downvote"
     ody: "*"

*/
func (s *cryptoVoteServiceServer) AddDownVote(ctx context.Context, req *pb.AddVoteReq) (*pb.TotalChangesRes, error) {
	var filterCryptoVote model.CryptoVote
	var err error

	// cria a variável de retorno
	var res *pb.TotalChangesRes = &pb.TotalChangesRes{
		Qtd: 0,
	}

	//retira do request um *pb.CryptoVote
	messageCryptoVote := req.Filter
	if messageCryptoVote != nil {
		filterCryptoVote, err = helpperPrepareFilter(messageCryptoVote)
		if err != nil {
			log.Printf("Problemas com helpperPrepareFilter: %v", err)
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			return res, status.Errorf(grpc_error_code, error_message)
		}
	} else {
		// se for nil precisamos retornar pois não pode loclaizar registro
		log.Printf("Problemas com o filtro de busca: %v", err)
		err = errors.New("[cryptovote.grpc_server] filtro não pode ser empty")
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}

	s.mutex.Lock()
	// cria um model.CryptoVote com dados da req
	cyptoVote, err := bo.AddDownVote(filterCryptoVote.Name, filterCryptoVote.Symbol)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		return res, status.Errorf(grpc_error_code, error_message)
	}
	s.mutex.Unlock()

	// se o ID não for zero sinal que conseguiu realizar DownVote
	if !cyptoVote.Id.IsZero() {
		res.Qtd = 1
	}
	return res, err
}

// loadCryptoVotesFromMongoDB carrega CryptoVotes
func (s *cryptoVoteServiceServer) loadCryptoVotes() {
	var modelCryptoVotes []model.CryptoVote

	// filtro vazio para recuperar todos os dados do banco
	var filterCryptoVote = model.CryptoVote{}

	modelCryptoVotes, err := bo.RetrieveAllCryptoVoteByFilter(filterCryptoVote)
	if err != nil {
		error_message := err.Error()
		grpc_error_code := status.Code(err)
		err = status.Errorf(grpc_error_code, error_message)
		log.Printf("[cryptovote.grpc_server] +%v", err)
	}
	if len(modelCryptoVotes) > 0 {

		// convertendo de []model.CryptoVote para json
		jsonData, err := json.Marshal(modelCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[cryptovote.grpc_server] +%v", err)
		}

		// convertendo de json para []pb.CryptoVote
		err = json.Unmarshal(jsonData, &s.cachedCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[cryptovote.grpc_server] +%v", err)
		}

	} else if len(jsonDefalutDataCryptoVotes) > 0 {
		// convertendo de json para []model.CryptoVote
		err := json.Unmarshal(jsonDefalutDataCryptoVotes, &modelCryptoVotes)
		if err != nil {
			error_message := err.Error()
			grpc_error_code := status.Code(err)
			err = status.Errorf(grpc_error_code, error_message)
			log.Printf("[cryptovote.grpc_server] +%v", err)
		}

		// depois do Unmarshal se tiver CryptoVote faz um create
		if len(modelCryptoVotes) > 0 {
			for _, modelCryptoVote := range modelCryptoVotes {

				modelCryptoVote, err = bo.CreateCryptoVote(modelCryptoVote)
				if err != nil {
					error_message := err.Error()
					grpc_error_code := status.Code(err)
					err = status.Errorf(grpc_error_code, error_message)
					log.Printf("[cryptovote.grpc_server] +%v", err)
				}
			}
		}
		//chama recursivamente para entrar no primeiro laço
		s.loadCryptoVotes()
	}
}

// imprime no console do servidor
func (s *cryptoVoteServiceServer) printServerInfo() {
	log.Printf("[cryptovote.grpc_server] Cryptovote service is a web server using gRPC")
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
		log.Printf("[cryptovote.grpc_server] failed to listen: %v", err)
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
			log.Printf("[cryptovote.grpc_server] Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	cryptovoteGrpcServer := grpc.NewServer(opts...)
	pb.RegisterCryptoVoteServiceServer(cryptovoteGrpcServer, newServer())

	// registranndo o servidor como reflection service para usar o gRPCuic
	reflection.Register(cryptovoteGrpcServer)

	log.Printf("[cryptovote.grpc_server] listening at %v", lis.Addr())

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
	"name": "bitcoin",
    "symbol": "btc",
	"qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"name": "ethereum",
    "symbol": "eth",
    "qtd_upvote": 0,
	"qtd_downvote": 0
}, {
	"name": "klever",
	"symbol": "klv",
    "qtd_upvote": 0,
	"qtd_downvote": 0	
}]`)
