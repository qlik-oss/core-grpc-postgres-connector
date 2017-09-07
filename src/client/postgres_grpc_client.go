package main

import (
	"../qlik"
	"google.golang.org/grpc"
	"context"
	"fmt"
	"time"
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
	var getDataOptions = &qlik.GetDataOptions{};


	var t0 = makeTimestamp();

	var stream, err2 = client.GetData2(context.Background(), getDataOptions)
	if err2 != nil {
		fmt.Println(err)
	}
	var bundle, receiveError = stream.Recv();
	if (bundle.CellsByRow != nil) {
		fmt.Println("cells by row");
	} else {
		fmt.Println("cells by column");
	}

	var totalCount int64;
	for receiveError == nil {
		if (bundle.CellsByRow != nil) {
			totalCount += bundle.CellsByRow.RowCount
		} else {
			totalCount += bundle.CellsByColumn.RowCount;
		}

		bundle, receiveError = stream.Recv()
	}
	var t1 = makeTimestamp();
	fmt.Println("Total rows", totalCount)
	fmt.Println("Time", t1 - t0, "ms")


}
