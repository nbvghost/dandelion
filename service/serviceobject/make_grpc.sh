export PATH=$PATH:../../third_software/protoc/
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative Grpc.proto