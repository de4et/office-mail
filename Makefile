
all: test

build-gateway:
	@echo "Building..."
	@go build -o ./bin/build-gateway ./services/mail-gateway/cmd/mail-gateway/main.go

run-gateway: build-gateway
	@./bin/build-gateway

test:
	@echo "Testing..."
	@go test ./... -v

SERVICES := \
	mail-gateway \
	mail-worker \
	delivery

docker-build-push:
	echo "Building and pushing images of all services..."
	echo $(SERVICES)
	for s in $(SERVICES); do \
		docker buildx build \
			--target $$s-production \
			-t de4et/$$s:latest \
			--push \
			.; \
	done

