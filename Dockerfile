FROM golang:1.13.8-alpine3.11

RUN mkdir -p go/src/github.com/dadosjusbr/alba/

WORKDIR /go/src/github.com/dadosjusbr/alba/

COPY . /go/src/github.com/dadosjusbr/alba/

RUN apk add bash

RUN apk add git

# TODO: Compilar o programa 
RUN go get go.mongodb.org/mongo-driver/mongo 
RUN go get github.com/urfave/cli
RUN go get github.com/urfave/cli/altsrc

RUN export GOBIN=$GOPATH/bin