GEN_GO = $$(which protoc-gen-go).exe
GEN_GO_GRPC = $$(which protoc-gen-go-grpc).exe

run-server:
	go run main.go server

gen:
	protoc -I=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO) --plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC) \
		./service/service.proto
