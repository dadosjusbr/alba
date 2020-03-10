FROM golang:1.13.8-alpine3.11

RUN apk add git

## TODO: compilar o programa aqui.

RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/urfave/cli
RUN go get github.com/urfave/cli/altsrc

RUN export GOBIN=$GOPATH/bin