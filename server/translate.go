package main

import (
	"encoding/json"
	"log"

	pb "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb"
	"github.com/caioformiga/go_mongodb_crud_cryptovote/model"
)

// TranslateToModelCryptoVote converte de *pb.CryptoVote para JSON e depois para model.CryptoVote
func TranslateToModelCryptoVote(messageCryptoVote *pb.CryptoVote) (model.CryptoVote, error) {
	var modelCryptoVote model.CryptoVote

	// convertendo de pb.CryptoVote para json
	jsonData, err := json.Marshal(messageCryptoVote)
	if err != nil {
		log.Fatalf("Problemas para fazer Marshal: %v", err)
	} else {

		// convertendo de json para model.CryptoVote
		err := json.Unmarshal(jsonData, &modelCryptoVote)
		if err != nil {
			log.Fatalf("Erro ao fazer Unmarshal de CryptoCurrencis em modelCryptoVote: %v", err)
		}
	}
	return modelCryptoVote, err
}

// TranslateToMessageCryptoVote de model.CryptoVote para JSON e depois para *pb.CryptoVote
func TranslateToMessageCryptoVote(modelCryptoVote model.CryptoVote) (*pb.CryptoVote, error) {
	var messageCryptoVote *pb.CryptoVote

	// convertendo de model.CryptoVote para json
	jsonData, err := json.Marshal(modelCryptoVote)
	if err != nil {
		log.Fatalf("Problemas para fazer Marshal: %v", err)
	} else {

		// convertendo de json para []pb.CryptoVote
		err := json.Unmarshal(jsonData, &messageCryptoVote)
		if err != nil {
			log.Fatalf("Erro ao fazer Unmarshal de CryptoCurrencis em messageDataCryptoCurrencis: %v", err)
		}
	}
	return messageCryptoVote, err
}
