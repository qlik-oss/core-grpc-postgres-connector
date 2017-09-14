FROM golang:1.9

ADD . /go/src/github.com/qlik-ea/postgres-grpc-connector

RUN go install github.com/qlik-ea/postgres-grpc-connector/server/

CMD /go/bin/server

EXPOSE 50051