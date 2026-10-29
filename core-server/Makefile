MODULE := core-server
PROTO_DIR := ./docs/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

wire:
	go run github.com/google/wire/cmd/wire ./cmd/server

generate-proto-rpc:
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(MODULE) \
		$(PROTO_FILES)

# protoc --proto_path=./docs/proto --go_out=. --go_opt=module=core-server --go-grpc_out=. --go-grpc_opt=module=core-server ./docs/proto/*.proto


ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/infras/repo/ent/schema
	go run -mod=mod entgo.io/ent/cmd/ent new User