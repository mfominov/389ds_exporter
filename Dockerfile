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

RUN CGO_ENABLED=0 go build -o /389ds_exporter

##
## Deploy
##
FROM golang:alpine

WORKDIR /

COPY --from=build /389ds_exporter /389ds_exporter

EXPOSE 9496

ENTRYPOINT ["/389ds_exporter"]
