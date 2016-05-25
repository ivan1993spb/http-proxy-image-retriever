FROM golang:1.6.2-wheezy

ADD . /go/src/github.com/ivan1993spb/http-proxy-image-retriever

RUN go install github.com/ivan1993spb/http-proxy-image-retriever

ENTRYPOINT /go/bin/http-proxy-image-retriever

EXPOSE 8888
