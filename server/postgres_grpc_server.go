package main

import (
	"fmt"
	"log"
	"github.com/qlik-ea/postgres-grpc-connector/qlik"
	"os"
	"runtime/pprof"
	"flag"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"github.com/qlik-ea/postgres-grpc-connector/postgres"
	"google.golang.org/grpc/codes"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")

type server struct {
	postgresReaders map[string]*postgres.PostgresReader
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (this *server) ExecuteGenericCommand(
	context context.Context, genericCommand *qlik.GenericCommand) (*qlik.GenericCommandResponse, error) {
	return &qlik.GenericCommandResponse{Data: "{}"}, nil
}

func (this *server) GetData(dataOptions *qlik.GetDataOptions, stream qlik.Connector_GetDataServer) error {
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
	var connectionString = dataOptions.Connection.ConnectionString

	if dataOptions.Connection.User != "" {
		connectionString = connectionString + ";user=" + dataOptions.Connection.User
	}

	if dataOptions.Connection.Password != "" {
		connectionString = connectionString + ";password=" + dataOptions.Connection.Password
	}

	if this.postgresReaders[connectionString] == nil {
		fmt.Println("Starting connection pool");
		fmt.Println(connectionString);
		reader, err := postgres.NewPostgresReader(connectionString)
		this.postgresReaders[connectionString] = reader
		if err != nil {
			err1 := status.Error(codes.Unavailable, "FOOO BAAAR")
			fmt.Println(err1)
			return err1
		}
	} else {
		fmt.Println("Reusing connection pool")
	}

	err := this.postgresReaders[connectionString].GetData(dataOptions, stream)
	t1 := makeTimestamp()
	fmt.Println("Time", t1-t0, "ms")

	return err
}
