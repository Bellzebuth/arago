PROTO_DIR=proto/tracker
OUT_DIR=proto/tracker

protoc:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/tracker.proto

