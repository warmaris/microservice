PROTO_INPUT_DIR := api/v1
PROTO_OUTPUT_DIR := pkg/v1
PROTO_MODULE_NAME := microservice

.PHONY: proto
proto:
	mkdir -p $(PROTO_OUTPUT_DIR)
	protoc \
		-I=api \
		--go_out=$(PROTO_OUTPUT_DIR) \
		--go_opt=module=$(PROTO_MODULE_NAME)/$(PROTO_OUTPUT_DIR) \
		--go-grpc_out=$(PROTO_OUTPUT_DIR) \
		--go-grpc_opt=module=$(PROTO_MODULE_NAME)/$(PROTO_OUTPUT_DIR) \
		$(shell ls -p $(PROTO_INPUT_DIR) | grep -v / | sed 's/^/v1\//')