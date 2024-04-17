FROM golang:1.20-alpine AS builder

WORKDIR /
COPY ./lb ./lb
COPY go.mod .
COPY go.sum .
COPY hosts.json .
RUN go build -o /balancer /lb/src/cmd/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /balancer /balancer
COPY --from=builder hosts.json hosts.json

EXPOSE 8080
CMD ["/balancer"]