gen:
	protoc --proto_path=proto --go_out=api --go_opt=paths=source_relative --go-grpc_out=api --go-grpc_opt=paths=source_relative data.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc --proto_path=proto --go_out=api --go_opt=paths=source_relative --go-grpc_out=api --go-grpc_opt=paths=source_relative auth_service.proto --go-grpc_opt=require_unimplemented_servers=false

clean:
	rm api/*.go

server:
	go run server/*.go

client:
	go run client/*.go

.PHONY: gen clean server client
