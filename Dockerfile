FROM golang as builder
COPY . /src
WORKDIR /src
ENV CGO_ENABLED 0
RUN go build -mod vendor -o /assets/in ./in
RUN go build -mod vendor -o /assets/out ./out
RUN go build -mod vendor -o /assets/check ./check
RUN set -e; \
      for pkg in $(go list -mod vendor ./...); do \
        go test -mod vendor -o "/tests/$(basename $pkg).test" -c $pkg; \
      done

FROM ubuntu:bionic AS resource
RUN apt-get update \
      && DEBIAN_FRONTEND=noninteractive \
      apt-get install -y --no-install-recommends \
        ca-certificates \
      && DEBIAN_FRONTEND=noninteractive \
        apt-get install -y git \
      && rm -rf /var/lib/apt/lists/*
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests /go-tests
WORKDIR /go-tests
RUN set -e; \
      for test in /go-tests/*.test; do \
        $test; \
      done

FROM resource
