# Images name: viettin/mai-admin:golang-base

FROM golang:1.15-alpine

ENV PATH=$PATH:$GOPATH/bin
RUN apk update && apk add --no-cache git ca-certificates

RUN go get -u github.com/swaggo/swag && \
    go get -u golang.org/x/crypto/bcrypt && \
    mkdir -p /go/src/github.com/D/gen/proto && \
    echo "Build proto done"

WORKDIR /Dcarbon
COPY . .
RUN go mod download