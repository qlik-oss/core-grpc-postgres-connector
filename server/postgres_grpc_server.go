package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/qlik-oss/core-grpc-postgres-connector/postgres"
	qlik "github.com/qlik-oss/core-grpc-postgres-connector/qlik"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")

type server struct {
	postgresReaders map[string]*postgres.Reader
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (s *server) GetData(dataRequest *qlik.DataRequest, stream qlik.Connector_GetDataServer) error {

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

	var t0 = makeTimestamp()
	var connectionString = dataRequest.Connection.ConnectionString

	if dataRequest.Connection.User != "" {
		connectionString = connectionString + ";user=" + dataRequest.Connection.User
	}

	if dataRequest.Connection.Password != "" {
		connectionString = connectionString + ";password=" + dataRequest.Connection.Password
	}

	if s.postgresReaders[connectionString] == nil {
		fmt.Println("Starting connection pool")
		fmt.Println(connectionString)
		reader, err := postgres.NewPostgresReader(connectionString)
		s.postgresReaders[connectionString] = reader
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	} else {
		fmt.Println("Reusing connection pool")
	}

	err := s.postgresReaders[connectionString].GetData(dataRequest, stream)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
	}
	t1 := makeTimestamp()
	fmt.Println("Time", t1-t0, "ms")

	return err
}

func (s *server) GetMetaInfo(ctx context.Context, metaInfoRequest *qlik.MetaInfoRequest) (*qlik.MetaInfo, error) {
	var metaInfo = qlik.MetaInfo{Name: "Postgres GPRC connector", Version: "1.0.0", Developer: "Qlik"}
	return &metaInfo, nil
}
