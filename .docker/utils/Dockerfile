ARG BUILD_ARG_GO_VERSION=1.23
ARG BUILD_ARG_ALPINE_VERSION=3.20

FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION} AS builder

ENV GOLANGCI_LINT_VERSION=1.60.1
ENV BUF_VERSION=1.30.1
ENV MOCKERY_VERSION=2.35.4
ENV GOWRAP_VERION=1.3.2
ENV PROTOC_GO_VERSION=1.31.0
ENV PROTOC_GRPC_VERSION=1.3.0
ENV GRPC_GATEWAY_VERSION=2.18.0
ENV GOFLAGS="-buildvcs=false"

RUN apk add --update --no-cache git gcc libc-dev ca-certificates tar

ADD https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz /tmp
ADD https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64 /tmp
ADD https://github.com/vektra/mockery/releases/download/v${MOCKERY_VERSION}/mockery_${MOCKERY_VERSION}_Linux_x86_64.tar.gz /tmp
ADD https://github.com/hexdigest/gowrap/releases/download/v${GOWRAP_VERION}/gowrap_${GOWRAP_VERION}_linux_amd64.tar.gz /tmp

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GO_VERSION} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GRPC_VERSION} && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION} && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v${GRPC_GATEWAY_VERSION}


RUN mkdir -p /tmp/utils && \
    tar -xzvf /tmp/golangci-lint-${GOLANGCI_LINT_VERSION}-linux-amd64.tar.gz -C /go/bin && \
    tar -xzvf /tmp/mockery_${MOCKERY_VERSION}_Linux_x86_64.tar.gz -C /go/bin && \
    tar -xzvf /tmp/gowrap_${GOWRAP_VERION}_linux_amd64.tar.gz -C /go/bin && \
    cp /tmp/buf-Linux-x86_64 /go/bin/buf && \
    chmod +x /go/bin/*

WORKDIR /src

FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION}

RUN apk add git
COPY --from=builder /go/bin/* /go/bin/

WORKDIR /src

RUN go env -w GOFLAGS=-buildvcs=false