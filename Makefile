.PHONY: build
	@echo "+ $@"
	go build -o bin/server cmd/main.go

.PHONY: proto
proto:
	protoc -I scripts/proto/ scripts/proto/api.proto --go_out=plugins=grpc:internal/api/