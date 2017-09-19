FROM golang:1.9

RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/qlik-ea/postgres-grpc-connector
WORKDIR /go/src/github.com/qlik-ea/postgres-grpc-connector
RUN dep ensure
RUN go install github.com/qlik-ea/postgres-grpc-connector/server
CMD ["/go/bin/server"]

EXPOSE 50051