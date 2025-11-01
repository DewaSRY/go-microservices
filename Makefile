PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

# POSTGRES_URI := postgres://postgres:postgres@localhost:5434/mydb?sslmode=disable
POSTGRES_URI := postgres://postgres:postgres@localhost:5433/riderdb?sslmode=disable

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

migration-up:
	goose -dir ./migrations postgres $(POSTGRES_URI) up

migration-down:
	goose -dir ./migrations postgres $(POSTGRES_URI) down