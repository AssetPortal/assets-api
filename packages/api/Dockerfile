FROM golang:1.23.3-alpine3.20 AS builder

COPY go.mod ./

RUN go mod download

COPY . /go/src/api

WORKDIR /go/src/api

RUN go build -o ./build/bin/api ./cmd/main.go

FROM alpine

ARG SERVER_PORT=8000

COPY --from=builder go/src/api/build/bin/api /usr/local/bin

WORKDIR /usr/local/bin

EXPOSE ${SERVER_PORT}

RUN apk --no-cache add ca-certificates

CMD ["api"]
