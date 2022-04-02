# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /389ds_exporter

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /389ds_exporter  /389ds_exporter

EXPOSE 9496

USER nonroot:nonroot

ENTRYPOINT ["/389ds_exporter"]
