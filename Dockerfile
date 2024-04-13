FROM golang:1.22.0-alpine3.19 AS builder

ENV GOPROXY=https://goproxy.io,direct
RUN apk update --no-cache
WORKDIR /usr/local/go/src
COPY . /usr/local/go/src
RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/s3/main.go

FROM alpine

RUN apk update --no-cache

WORKDIR /app

COPY --from=builder /usr/local/go/src/app /app

CMD ./app