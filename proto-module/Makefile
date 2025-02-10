PROTO_DIR=proto

SERVICE=auth
# SERVICE=payment
# SERVICE=order
# SERVICE=notifcation
# SERVICE=product
# SERVICE=user

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/$(SERVICE)/$(SERVICE).proto
