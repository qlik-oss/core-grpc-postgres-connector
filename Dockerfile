FROM golang:1.9
RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/github.com/qlik-ea/postgres-grpc-connector
WORKDIR /go/src/github.com/qlik-ea/postgres-grpc-connector
RUN dep ensure
# RUN go install github.com/qlik-ea/postgres-grpc-connector/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./server 

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/github.com/qlik-ea/postgres-grpc-connector/main .
CMD ["./main"]
EXPOSE 50051