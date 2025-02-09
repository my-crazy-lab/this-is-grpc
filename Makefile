PROTO_DIR=proto
OUT_DIR=graphql-api-gateway
# OUT_DIR=authentication

generate:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	       $(PROTO_DIR)/account/account.proto
