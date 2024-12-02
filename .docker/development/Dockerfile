FROM golang:1.23 as base
RUN apt install git gcc libc-dev ca-certificates make iostat

WORKDIR /src
COPY . .

RUN go get -d -t

FROM base as builder
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-w -s" -o build/release/ .

FROM base as dev
WORKDIR /app
EXPOSE 8081
CMD ["tail", "-f", "/dev/null"]

FROM ubuntu:18.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends\
        ca-certificates \
        tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN ln -snf /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
    && echo Europe/Moscow > /etc/timezone
RUN addgroup -gid 8000 app \
    && useradd -g app app -u 8000 -m -d /app -s /bin/bash

WORKDIR /app
COPY --from=builder /src/build/release .

USER app:app