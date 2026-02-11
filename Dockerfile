FROM golang:1.26.0

WORKDIR /go/src

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin v2.5.0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD ["make", "run"]
