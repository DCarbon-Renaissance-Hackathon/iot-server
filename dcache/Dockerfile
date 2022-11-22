FROM golang:1.19-alpine as builder

WORKDIR /dcarbon
COPY . .
RUN apk add --no-cache alpine-sdk

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go mod download 

ENV PATH=$PATH:$GOPATH/bin

