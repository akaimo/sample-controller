FROM golang:1.12.9 as builder

ENV GO111MODULE=on
WORKDIR /go/src/github.com/akaimo.com/sample-controller

COPY go.mod .
COPY go.sum .

RUN set -x \
  && go mod download

COPY . .

RUN go build -o sample-controller .

FROM ubuntu:bionic

WORKDIR /
COPY --from=builder /go/src/github.com/akaimo.com/sample-controller/sample-controller .

CMD ["/sample-controller"]
