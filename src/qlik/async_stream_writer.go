package qlik

/**
 *	Class AsyncStreamWriter
 */
type AsyncStreamWriter struct {
	grpcStream Connector_GetData2Server
	channel chan *ResultChunk
	done *chan bool
}

func NewAsyncStreamWriter(grpcStream Connector_GetData2Server, done *chan bool) *AsyncStreamWriter {
	var this = &AsyncStreamWriter{grpcStream, make(chan *ResultChunk, 100), done}
	go this.run();
	return this
}
func (this *AsyncStreamWriter) Write(rowBundle *ResultChunk) {
	this.channel <- rowBundle
}

func (this *AsyncStreamWriter) Close() {
	close(this.channel)
}

func (this *AsyncStreamWriter) run() {
	for resultChunk := range this.channel {
		this.grpcStream.Send(resultChunk);
	}
	*this.done <- true;
}

