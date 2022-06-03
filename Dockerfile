FROM golang:1.18.3

WORKDIR /go/src

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin v1.45.2

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD make run
