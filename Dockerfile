# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.16 AS build

ARG GIT_COMMIT

WORKDIR /go/src/github.com/figment-networks/avalanche-indexer

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV GOARCH=amd64
ENV GOOS=linux

RUN make setup

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$GIT_COMMIT \
  make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.14 AS release

WORKDIR /app

COPY --from=build \
  /go/src/github.com/figment-networks/avalanche-indexer/avalanche-indexer \
  /app/avalanche-indexer

EXPOSE 8081

ENTRYPOINT ["/app/avalanche-indexer"]