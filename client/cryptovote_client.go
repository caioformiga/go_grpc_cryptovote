package main

import (
	"context"
	"flag"
	"io"

	//"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	pb "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

// printFeature gets the feature for the given point.
func printAllCrypto(client pb.CryptoVoteServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Faz uma chamada a função rpc ListAllCryptoCurrencies detalhes em go_grpc_cryptovote_grpc.pb.go ")
	stream, err := client.ListAllCryptoCurrencies(ctx, &pb.EmptyReq{})
	log.Printf("Finanlizou!")

	log.Printf("Recuperando todoas as Crypto's...")
	if err != nil {
		log.Fatalf("%v.ListAllCryptoCurrencies(_) = _, %v", client, err)
	}
	for {
		crypto, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListAllCryptoCurrencies(_) = _, %v", client, err)
		}
		log.Printf("Crypto Name:%v\n", crypto.GetName())
		log.Printf("Crypto Symbol:%v\n", crypto.GetSymbol())
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())

	// tenta criar uma conexao com o referido serverAddr = "localhost:10000"
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Nao conseguiu fazer conexao com o servidor: %v", err)
	}
	defer conn.Close()
	cryptovoteGrpcClient := pb.NewCryptoVoteServiceClient(conn)

	printAllCrypto(cryptovoteGrpcClient)
}
