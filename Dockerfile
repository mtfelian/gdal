FROM golang:1.18.7 as builder
RUN apt update && apt install -y libgdal-dev
WORKDIR '/app'
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build rm go.sum && go mod tidy && go build
