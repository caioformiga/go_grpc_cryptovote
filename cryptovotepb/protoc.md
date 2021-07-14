The .proto file exposes a CryptoVoteService which features four rpc methods

    // Obtains all CryptoCurrency    
    rpc ListAllCryptoCurrencies(EmptyReq) returns (GetCryptocurrencyResponse) {}

    rpc GetStreamCryptoCurrencies(EmptyReq) returns (stream Cryptocurrency) {}


    // Obtains all CryptoVote
    rpc ListAllCryptoVotes(EmptyReq) returns (GetCryptoVoteResponse) {}

    rpc GetStreamCryptoVotes(EmptyReq) returns (stream CryptoVote) {}


This methods can be called by any gRPC client written in any language.
These .proto definitions are typically shared across clients of all 
shapes and sizes so that they can generate their own code to talk to 
our gRPC server.


gRPC uses protocol buffer language, for each source .proto file must 
have a single Go package, witch is specified "go_package" at .proto file

    option go_package = "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb";

There is a one-to-one relationship between source .proto files and .pb.go files,
protoc compiler replace the .proto suffix with .pb.go
    go_grpc_cryptovote.proto

    go_grpc_cryptovote.pb.go


However, the output directory can be selected:

    cryptovotepb/go_grpc_cryptovote.proto 
    option go_package = "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb";
    
The corresponding output file may be two ways (Relative to the import path or Relative to the input file)

Use the second, Relative to the input file:
command: protoc --go_out=paths=source_relative:. cryptovotepb/go_grpc_cryptovote.proto
result: $GOPATH/src/github.com/caioformiga/go_grpc_cryptovote/cryptovotepb/go_grpc_cryptovote.pb.go

#1 Relative to the import path:
command: protoc --go_out=. cryptovotepb/go_grpc_cryptovote.proto
result: $GOPATH/src/github.com/caioformiga/go_grpc_cryptovote/github.com/caioformiga/go_grpc_cryptovote/cryptovotepb/go_grpc_cryptovote.pb.go

link: https://chromium.googlesource.com/external/github.com/golang/protobuf/+/refs/tags/v1.1.0/README.md


cryptovotepb/go_grpc_cryptovote.proto 
option go_package = "github.com/caioformiga/go_grpc_cryptovote/cryptovotepb";

cd $GOPATH/src/github.com/caioformiga/go_grpc_cryptovote

protoc --go_out=paths=source_relative:. cryptovotepb/go_grpc_cryptovote.proto
result: $GOPATH/src/github.com/caioformiga/go_grpc_cryptovote/cryptovotepb/go_grpc_cryptovote.pb.go

protoc --go-grpc_out=paths=source_relative:. cryptovotepb/go_grpc_cryptovote.proto
result: $GOPATH/src/github.com/caioformiga/go_grpc_cryptovote/cryptovotepb/go_grpc_cryptovote_grpc.pb.go
