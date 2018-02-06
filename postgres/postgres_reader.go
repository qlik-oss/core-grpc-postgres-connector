package postgres

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/jackc/pgx"
	"github.com/qlik-ea/postgres-grpc-connector/qlik"
	"google.golang.org/grpc/metadata"
)

// Reader contains the pool of postegres connections.
type Reader struct {
	pool *pgx.ConnPool
}

// NewPostgresReader constructs a new reader.
func NewPostgresReader(connectString string) (*Reader, error) {
	var pool, err = pgx.NewConnPool(extractConfig(connectString))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		return nil, err
	}
	return &Reader{pool}, nil
}

// GetData will return data from the postgres database.
func (r *Reader) GetData(dataOptions *qlik.GetDataOptions, stream qlik.Connector_GetDataServer) error {
	var done = make(chan bool)
	// Connect to postgres
	conn, err := r.pool.Acquire()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
	}
	defer r.pool.Release(conn)

	// Select postgresRowData

	fmt.Println(dataOptions.Parameters.Statement)
	fmt.Println(dataOptions.Connection.ConnectionString)
	fmt.Println(dataOptions.Connection.User)
	rows, err2 := conn.Query(dataOptions.Parameters.Statement)
	if err2 != nil {
		fmt.Println(err2)
	}

	// Start asynchronus translation and writing
	var asyncStreamwriter = qlik.NewAsyncStreamWriter(stream, done)
	var asyncTranslator = NewAsyncTranslator(asyncStreamwriter, rows.FieldDescriptions())

	// Set header with postgresRowData format
	var headerMap = make(map[string]string)
	var getDataResponseBytes, _ = proto.Marshal(asyncTranslator.GetDataResponseMetadata())
	headerMap["x-qlik-getdata-bin"] = string(getDataResponseBytes)
	stream.SendHeader(metadata.New(headerMap))

	//Read postgresRowData from postgres
	const MaxRowsPerBundle = 200
	var rowList = [][]interface{}{}
	for rows.Next() {
		var srcColumns, _ = rows.Values()
		rowList = append(rowList, srcColumns)
		if len(rowList) >= MaxRowsPerBundle {
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
	return nil
}

func extractConfig(connectString string) pgx.ConnPoolConfig {
	params := connectStringToParamsMap(connectString)
	var config pgx.ConnPoolConfig

	config.Host = params["host"]
	if config.Host == "" {
		config.Host = params["hostname"]
	}
	if params["port"] != "" {
		var intPort, _ = strconv.Atoi(params["port"])
		config.Port = uint16(intPort)
	}

	config.User = params["username"]
	if config.User == "" {
		config.User = params["user"]
	}
	if config.User == "" {
		config.User = params["userid"]
	}

	config.Password = params["password"]
	config.Database = params["database"]

	return config
}
func connectStringToParamsMap(connectString string) map[string]string {
	var params = strings.Split(connectString, ";")
	paramsMap := make(map[string]string)
	for _, v := range params {
		paramAndValue := strings.Split(v, "=")
		if len(paramAndValue) == 2 {
			param := strings.ToLower(strings.TrimSpace(paramAndValue[0]))
			value := strings.TrimSpace(paramAndValue[1])
			paramsMap[param] = value
		}
	}
	return paramsMap
}
