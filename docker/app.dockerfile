FROM golang:alpine

WORKDIR /Go-MCS

ADD . .

RUN go mod download

ENTRYPOINT go build  && ./Go-MCS