FROM golang:1.22.1-alpine as builder

COPY . /app/
WORKDIR /app/

RUN go mod download
RUN go build -o ./bin/server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/bin/server .

CMD ["./server"]