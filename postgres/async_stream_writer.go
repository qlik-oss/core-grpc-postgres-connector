package postgres

import (
	qlik "github.com/qlik-ea/postgres-grpc-connector/qlik"
)

// AsyncStreamWriter defines the writer interface.
type AsyncStreamWriter struct {
	grpcStream qlik.Connector_GetDataServer
	channel    chan *qlik.DataChunk
	done       chan bool
}

// NewAsyncStreamWriter constructs a new async stream writer.
func NewAsyncStreamWriter(grpcStream qlik.Connector_GetDataServer, done chan bool) *AsyncStreamWriter {
	var this = &AsyncStreamWriter{grpcStream, make(chan *qlik.DataChunk, 10), done}
	go this.run()
	return this
}

// Write will send a datachunk on the underlying stream.
func (a *AsyncStreamWriter) Write(rowBundle *qlik.DataChunk) {
	a.channel <- rowBundle
}

// Close will close the underlying stream.
func (a *AsyncStreamWriter) Close() {
	close(a.channel)
}

func (a *AsyncStreamWriter) run() {
	for resultChunk := range a.channel {
		a.grpcStream.Send(resultChunk)
	}
	a.done <- true
}
