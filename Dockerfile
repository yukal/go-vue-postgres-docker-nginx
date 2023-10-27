# Documentation:
# https://docs.docker.com/engine/reference/builder/
#
# docker build -t <name>/<subname>:<tag> .
# docker build -t yu/goweb:v1 .

FROM golang:alpine AS build

RUN mkdir /src
WORKDIR /src

COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum ./
RUN go mod download && go mod verify \
  && go build ./cmd/server \
  && rm -R cmd internal go.*

FROM alpine:edge

WORKDIR /usr/sbin
COPY --from=build /src/server ./api

EXPOSE 50598

CMD /usr/sbin/api
