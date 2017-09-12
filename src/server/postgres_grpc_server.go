/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"fmt"
	"log"
	"net"
	"../qlik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/jackc/pgx"
	"os"
	"runtime/pprof"
	"flag"
	"time"
	"google.golang.org/grpc/metadata"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
)

const (
	port = ":50051"
)
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
var pool *pgx.ConnPool
type server struct{}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}





func (s *server) ExecuteGenericCommand(context context.Context, genericCommand *qlik.GenericCommand) (*qlik.GenericCommandResponse, error) {
	return &qlik.GenericCommandResponse{Data: "{}"}, nil
}

func (s *server) GetData(dataOptions *qlik.GetDataOptions, stream qlik.Connector_GetDataServer) error {
	var done = make(chan bool)

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		fmt.Println("Started cpu profiling...")
		defer pprof.StopCPUProfile()
	}

	// Connect to postgres
	conn, err := pool.Acquire()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
	}
	defer pool.Release(conn)

	// Select data

	fmt.Println(dataOptions.Parameters.Statement);
	rows, _ := conn.Query(dataOptions.Parameters.Statement)

	var t0 = makeTimestamp()

	// Start asynchronus translation and writing
	var asyncStreamwriter = qlik.NewAsyncStreamWriter(stream, &done)
	var asyncTranslator = qlik.NewAsyncTranslator(asyncStreamwriter, rows.FieldDescriptions())

	// Set header with data format
	var headerMap = make(map[string]string)
	var getDataResponseBytes, _ = proto.Marshal(asyncTranslator.GetDataResponseMetadata());
	headerMap["x-qlik-getdata-bin"] = string(getDataResponseBytes)
	stream.SendHeader(metadata.New(headerMap))

	//Read data from postgres
	const MAX_ROWS_PER_BUNDLE = 200
	var rowList = [][]interface{}{}
	for rows.Next() {
		var srcColumns, _ = rows.Values()
		rowList = append(rowList, srcColumns)
		if len(rowList) >= MAX_ROWS_PER_BUNDLE {
			asyncTranslator.Write(rowList)
			rowList = [][]interface{}{}
		}
	}
	if len(rowList) > 0 {
		asyncTranslator.Write(rowList)
		rowList = [][]interface{}{}
	}
	asyncTranslator.Close()

	//Wait for all translater and writer to finish
	<-done
	var t1 = makeTimestamp()
	fmt.Println("Time", t1 - t0, "ms")
	return nil
}



func main() {
	fmt.Println("Started...")
	var err error
	pool, err = pgx.NewConnPool(extractConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	qlik.RegisterConnectorServer(s, &server{})
	fmt.Println("Server registered...")
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//PGHOST=localhost;PGUSER=testuser;PGPASSWORD=testuser;PGDATABASE=test
func extractConfig() pgx.ConnPoolConfig {
	var config pgx.ConnPoolConfig

	config.Host = os.Getenv("PGHOST")
	if config.Host == "" {
		config.Host = "localhost"
	}

	config.User = os.Getenv("PGUSER")
	if config.User == "" {
		config.User = "testuser"
	}

	config.Password = os.Getenv("PGPASSWORD")
	if config.Password == "" {
		config.Password = "testuser"
	}

	config.Database = os.Getenv("PGDATABASE")
	if config.Database == "" {
		config.Database = "test"
	}

	return config
}
