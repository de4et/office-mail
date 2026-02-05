FROM golang:1.25-alpine AS go-builder

WORKDIR /app

COPY go.mod go.sum .

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY pkg/ pkg/
COPY services/mail-gateway/ ./services/mail-gateway/

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" \
    -o /mail-gateway ./services/mail-gateway/cmd/mail-gateway/main.go

COPY services/mail-worker/ ./services/mail-worker/

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" \
    -o /mail-worker ./services/mail-worker/cmd/mail-worker/main.go

FROM scratch AS mail-gateway-production

COPY .env .env
COPY --from=go-builder /mail-gateway /mail-gateway

ENTRYPOINT ["/mail-gateway"]

FROM scratch AS mail-worker-production

COPY .env .env
COPY --from=go-builder /mail-worker /mail-worker

ENTRYPOINT ["/mail-worker"]
