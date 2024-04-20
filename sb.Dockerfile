FROM golang:1.21-alpine as BUILDER
WORKDIR /

COPY ./s_server ./s_server
COPY go.mod .
COPY go.sum .

RUN go build -o /stupid_s /s_server/src/*.go



FROM alpine

WORKDIR /app
COPY --from=BUILDER /stupid_s /stupid_s

EXPOSE 8000
CMD ["/stupid_s"]


