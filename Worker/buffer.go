package Worker

import (
	"fmt"
	"io"
)

type FileContent struct {
	offset int64
	buffer io.ReadCloser
}

type Buffer struct {
	size   int
	buffer chan FileContent
	writer *FileWriter
}

func (this *Buffer) Init(size int, writer *FileWriter) {
	this.size = size
	this.buffer = make(chan FileContent, size)
	this.writer = writer
}

func (this *Buffer) Done(work FileContent) {
	if len(this.buffer) == this.size {
		fmt.Println("Done flush", len(this.buffer), this.size)
		this.Flush()
	}
	this.buffer <- work
}

func (this *Buffer) Flush() {
	for i := range this.buffer {
		this.writer.Write(i.offset, i.buffer)
	}
}

func (this *Buffer) Close() {
	close(this.buffer)
	if len(this.buffer) > 0 {
		this.Flush()
	}
	this.writer.Close()
}
