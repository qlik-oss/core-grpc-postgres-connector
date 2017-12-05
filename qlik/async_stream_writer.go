package qlik

// AsyncStreamWriter provides an API to asynchronously write data to a stream.
type AsyncStreamWriter struct {
	grpcStream Connector_GetDataServer
	channel    chan *DataChunk
	done       chan bool
}

// NewAsyncStreamWriter constructs a new async stream writer.
func NewAsyncStreamWriter(grpcStream Connector_GetDataServer, done chan bool) *AsyncStreamWriter {
	var this = &AsyncStreamWriter{grpcStream, make(chan *DataChunk, 10), done}
	go this.run()
	return this
}

// Write will send a datachunk on the underlying stream.
func (a *AsyncStreamWriter) Write(rowBundle *DataChunk) {
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
