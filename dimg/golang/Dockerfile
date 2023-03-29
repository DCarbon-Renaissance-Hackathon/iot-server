# Images name: viettin/mai-admin:golang-base

FROM golang:1.19-alpine

RUN apk update && apk add --no-cache git ca-certificates alpine-sdk

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    echo "Build proto done"

WORKDIR /dcarbon
COPY . .
RUN go mod download

ENV PATH=$PATH:$GOPATH/bin