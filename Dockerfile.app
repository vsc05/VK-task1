FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY app/main.go .

RUN go build -o /app/hello-world main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/hello-world .

CMD ["./hello-world"]
