FROM golang:1.13.8-alpine3.11

WORKDIR /go/src/github.com/dadosjusbr/alba/

COPY . .

RUN apk add git

RUN cd storage && go install

# TODO criar pacote da cli
RUN go get github.com/urfave/cli
RUN go get github.com/urfave/cli/altsrc

RUN export GOBIN=$GOPATH/bin