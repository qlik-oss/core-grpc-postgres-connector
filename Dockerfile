# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.9

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/qlik-trial/postgres-grpc-connector

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN /go/src/github.com/qlik-trial/postgres-grpc-connector/install-deps.sh

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/jackc/pgx
RUN go install github.com/qlik-trial/postgres-grpc-connector/server/

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/server

# Document that the service listens on port 8080.
EXPOSE 50051