package main

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	qlik "github.com/qlik-ea/postgres-grpc-connector/qlik"
	"google.golang.org/grpc"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	client := qlik.NewConnectorClient(conn)
	var dataRequest = &qlik.DataRequest{}
	dataRequest.Connection = &qlik.ConnectionInfo{
		ConnectionString: "host=selun-gwe.qliktech.com;database=test",
		User:             "testuser",
		Password:         "testuser",
	}
	dataRequest.Parameters = &qlik.DataInfo{
		Statement:  "select * from manytypes",
		Parameters: []*qlik.Parameter{},
	}
	var t0 = makeTimestamp()

	var stream, err2 = client.GetData(context.Background(), dataRequest)
	fmt.Println(err2)
	var header, err3 = stream.Header()
	fmt.Println(err3)
	var t = header["x-qlik-getdata-bin"]
	var t2 = t[0]
	var dataResponse = qlik.GetDataResponse{FieldInfo: make([]*qlik.FieldInfo, 100), TableName: "x"}
	proto.Unmarshal([]byte(t2), &dataResponse)
	fmt.Println("a", t)

	if err2 != nil {
		fmt.Println(err)
	}
	var dataChunk, receiveError = stream.Recv()

	var totalCount int
	for receiveError == nil {
		totalCount += len(dataChunk.StringCodes)

		dataChunk, receiveError = stream.Recv()
	}
	var t1 = makeTimestamp()
	fmt.Println("Total entries", totalCount)
	fmt.Println("Time", t1-t0, "ms")

}
