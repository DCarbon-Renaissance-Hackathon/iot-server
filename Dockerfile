FROM harbor.viet-tin.com/dcarbon/go-shared as builder

WORKDIR /dcarbon/iott-cloud
COPY . .

RUN swag init -g ./cmd/iott-cloud/main.go -o ./cmd/iott-cloud/docs  &&  \
    cd ./cmd/iott-cloud/ && \
    go mod tidy && \
    go build -buildvcs=false -o iott-cloud && \
    cp  iott-cloud /usr/bin


FROM harbor.viet-tin.com/dcarbon/dimg:minimal

COPY --from=builder /usr/bin/iott-cloud /usr/bin/iott-cloud
ENV GIN_MODE=release

CMD [ "iott-cloud" ]