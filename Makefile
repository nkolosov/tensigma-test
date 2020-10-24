OUTPUT?=bin/products-api

.PHONY: clean
clean:
	@echo "+ $@"
	rm -rf ${OUTPUT}

.PHONY: build
build: clean
	@echo "+ $@"
	go build -o bin/products-api cmd/server/server.go

.PHONY: proto
proto:
	@echo "+ $@"
	protoc --proto_path=./proto --go_out=plugins=grpc:./internal/api products.proto
