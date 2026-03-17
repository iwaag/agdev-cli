FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/agdev .

FROM alpine:3.20 AS cli

RUN adduser -D -u 10001 appuser
USER appuser

COPY --from=builder /out/agdev /usr/local/bin/agdev

ENTRYPOINT ["/usr/local/bin/agdev"]
