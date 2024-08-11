# Non alpine because of CGO glibc dependency
FROM golang:1.22 AS builder

WORKDIR /build
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -o ./openchangelog cmd/server.go

FROM alpine

ARG config=i-should-never-exists.jla
# Try to copy config, the * makes sure we don't fail if the file isn't found
COPY *$config /etc/openchangelog.yaml

WORKDIR /app
COPY --from=builder /build/openchangelog ./openchangelog

ENTRYPOINT ["/app/openchangelog"]