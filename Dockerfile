FROM golang:1.27.0

WORKDIR /go/src

# renovate: datasource=github-releases depName=golangci/golangci-lint
ARG GOLANGCI_LINT_VERSION=v2.13.1
RUN curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

CMD ["make", "run"]
