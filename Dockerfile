FROM golang:1.13.8-alpine3.11

RUN apk add git

RUN go get go.mongodb.org/mongo-driver/mongo