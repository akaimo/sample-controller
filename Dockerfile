FROM golang:1.13.3 as builder

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
