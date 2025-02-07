include .env
LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-shortener-api

generate-shortener-api:
	mkdir -p pkg/shortener
	protoc --proto_path api/shortener \
	--go_out=pkg/shortener --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/shortener --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/shortener/shortener.proto

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v