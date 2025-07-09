# Multi-platform build with CGO support
FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

# Install ARM64 cross-compilation toolchain
RUN apt-get update && apt-get install -y \
    gcc-aarch64-linux-gnu \
    libc6-dev-arm64-cross \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

# Build binary for target architecture
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

RUN if [ "$TARGETARCH" = "arm64" ]; then \
      export CC=aarch64-linux-gnu-gcc && \
      export CGO_CFLAGS="-g -O2" && \
      export CGO_CXXFLAGS="-g -O2" && \
      export CGO_LDFLAGS="-g -O2" && \
      go build -buildvcs=false -ldflags "-s -w" -o ./openchangelog cmd/server.go; \
    else \
      go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -o ./openchangelog cmd/server.go; \
    fi

FROM alpine

ARG config=i-should-never-exists.jla
# Try to copy config, the * makes sure we don't fail if the file isn't found
COPY *$config /etc/openchangelog.yaml

WORKDIR /app
COPY --from=builder /build/openchangelog ./openchangelog

ENTRYPOINT ["/app/openchangelog"]