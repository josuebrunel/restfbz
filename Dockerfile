FROM golang:1.17

ENV GO111MODULE=on

RUN     mkdir /go/src/app
WORKDIR /go/src/app
COPY    go.mod /app
ADD     . /go/src/app
# EXPOSE 8999
