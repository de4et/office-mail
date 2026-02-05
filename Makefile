
all: test

build-gateway:
	@echo "Building..."
	@go build -o ./bin/build-gateway ./services/mail-gateway/cmd/mail-gateway/main.go

run-gateway: build-gateway
	@./bin/build-gateway

test:
	@echo "Testing..."
	@go test ./... -v
