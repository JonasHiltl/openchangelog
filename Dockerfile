# https://github.com/tonistiigi/xx helpers for cross compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# copy helper scripts
COPY --from=xx / /
ARG TARGETPLATFORM

RUN apk add clang lld
RUN xx-apk add musl-dev gcc

WORKDIR /build
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

RUN xx-go build -buildvcs=false -ldflags "-s -w" -o ./openchangelog cmd/server.go && \
  xx-verify openchangelog

FROM alpine

ARG config=i-should-never-exists.jla
# Try to copy config, the * makes sure we don't fail if the file isn't found
COPY *$config /etc/openchangelog.yaml

WORKDIR /app
COPY --from=builder /build/openchangelog ./openchangelog

ENTRYPOINT ["/app/openchangelog"]