FROM golang:1.18.2-alpine AS builder

# Install all required dependencies for building
RUN apk update
RUN apk add git build-base
WORKDIR /service

## Install go tools
RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN task compile-restSimpleApp

FROM alpine:3.14.2
COPY --from=builder ./service/build/out/restSimpleApp restSimpleApp
ENTRYPOINT ["./restSimpleApp"]