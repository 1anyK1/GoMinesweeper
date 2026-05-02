FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server
RUN go build -o client ./cmd/client


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/client /app/client

ENV TCP_ADDR=:8080
ENV GAME_SIZE=8
ENV MINES=10

EXPOSE 8080

CMD ["/app/server"]