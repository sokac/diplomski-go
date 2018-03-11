FROM golang

ADD . /go/src/fibonacci

RUN cd /go/src/fibonacci && go install .

ENTRYPOINT ["/go/bin/fibonacci"]
