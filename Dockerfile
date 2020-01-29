FROM golang:alpine AS builder
EXPOSE 8080

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/pkg/app/
COPY . .

WORKDIR $GOPATH/src/pkg/app
RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine

RUN apk update && apk add --no-cache ca-certificates
EXPOSE 8080

WORKDIR /go/bin
ADD https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.9.6/kubeseal-linux-amd64 /usr/local/bin/kubeseal
RUN chmod +x /usr/local/bin/kubeseal && \
kubeseal --version
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /go/src/pkg/app/swaggerui /go/bin/swaggerui

ENTRYPOINT ["./app"]