# Builder stage
FROM golang:1.26-alpine AS builder

WORKDIR /builder/app

COPY ./go.mod ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o stress-test ./cmd/stress-test

# Final stage
FROM alpine:3.23
WORKDIR /app

# certificados p/ conseguir testar URLs https
RUN apk add --no-cache ca-certificates

COPY --from=builder /builder/app/stress-test .

# ENTRYPOINT (não CMD) para que as flags do `docker run` sejam
# repassadas ao binário: docker run imagem --url=... --requests=...
ENTRYPOINT ["./stress-test"]
