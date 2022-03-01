FROM golang:1.17-alpine as builder

WORKDIR /dcarbon
COPY . .
RUN apk add --no-cache alpine-sdk
RUN go build -o iott-cloud && cp  iott-cloud /usr/bin


FROM alpine:3.14

COPY --from=builder /usr/bin/iott-cloud /usr/bin/iott-cloud
ENV GIN_MODE=release

CMD [ "iott-cloud" ]