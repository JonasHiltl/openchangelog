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

WORKDIR /app
COPY --from=builder /build/openchangelog ./openchangelog

ENTRYPOINT ["/app/openchangelog"]