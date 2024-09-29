FROM golang:1.23.1 as build
LABEL builder=builder
WORKDIR /build
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download && \
    go mod verify
COPY . .
# go test -v ./...
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -o \
    /usr/bin/effectiveSong ./cmd/server/server.go

FROM alpine:3.12
RUN apk update && \
    apk add --no-cache \
        ca-certificates \
        curl \
        tzdata \
    && rm -rf -- /var/cache/apk/*
ENV TZ="UTC"
WORKDIR /app
COPY --from=build /usr/bin/effectiveSong .
COPY migrations ./migrations/
COPY docs/swagger.json ./docs/swagger.json
COPY config.override.json ./config.override.json
HEALTHCHECK --interval=20s --timeout=5s --retries=4 --start-period=20s \
    CMD curl -fsS -m5 -A'docker-healthcheck' http://127.0.0.1/api/v1/ping
CMD ["./effectiveSong"]

