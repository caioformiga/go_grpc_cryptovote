package main

import (
	"log"
	"testing"
)

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

func TestRetrieveAllCryptoVoteByFilter(t *testing.T) {
	//RetrieveAllCryptoVoteByFilter()
	resultado := "f"
	esperado := "Olá, mundo"

	if resultado != esperado {
		t.Errorf("resultado '%s', esperado '%s'", resultado, esperado)
	}
}
