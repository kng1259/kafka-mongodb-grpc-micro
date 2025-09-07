protodir=../proto
outdir=./protogen

protoc --proto_path=$protodir --go_out=./$outdir --go_opt=paths=source_relative --go-grpc_out=./$outdir --go-grpc_opt=paths=source_relative $protodir/product.proto