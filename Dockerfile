FROM harbor.viet-tin.com/dcarbon/iott-cloud:cache as builder

WORKDIR /dcarbon/iott-cloud
COPY . .

RUN apk add --no-cache alpine-sdk

RUN swag init -g ./cmd/iott-cloud/main.go -o ./cmd/iott-cloud/docs  &&  \
    cd ./cmd/iott-cloud/ && \
    go mod tidy && \
    go build -buildvcs=false -o iott-cloud && \
    cp  iott-cloud /usr/bin


FROM alpine:3.16

COPY --from=builder /usr/bin/iott-cloud /usr/bin/iott-cloud
ENV GIN_MODE=release

CMD [ "iott-cloud" ]