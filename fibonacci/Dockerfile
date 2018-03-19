FROM golang

ADD . /go/src/fibonacci

RUN cd /go/src/fibonacci && go build .

WORKDIR /go/src/fibonacci

ENTRYPOINT ["/go/src/fibonacci/fibonacci"]
