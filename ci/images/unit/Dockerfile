FROM golang:1.11-alpine
RUN true \
  && apk --no-cache add bash ca-certificates curl git openssl \
  && go get github.com/mattn/goveralls \
  && go get github.com/modocache/gover \
  && go get github.com/onsi/ginkgo/ginkgo
