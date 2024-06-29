PROTO_DIR = proto
PROTO_SRC = $(PROTO_DIR)/auth/owners.proto
GEN_DIR = gen/go

all: generate

generate:
	protoc -I $(PROTO_DIR) $(PROTO_SRC) --go_out=./$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=./$(GEN_DIR)/ --go-grpc_opt=paths=source_relative

clean:
	rm -rf $(GEN_DIR)
