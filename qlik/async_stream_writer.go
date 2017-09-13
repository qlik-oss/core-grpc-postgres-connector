package qlik

/**
 *	Class AsyncStreamWriter
 */
type AsyncStreamWriter struct {
	grpcStream Connector_GetDataServer
	channel chan *BundledRows
	done *chan bool
}

func NewAsyncStreamWriter(grpcStream Connector_GetDataServer, done *chan bool) *AsyncStreamWriter {
	var this = &AsyncStreamWriter{grpcStream, make(chan *BundledRows, 10), done}
	go this.run();
	return this
}
func (this *AsyncStreamWriter) Write(rowBundle *BundledRows) {
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

