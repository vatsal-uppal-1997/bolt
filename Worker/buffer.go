package Worker

import (
	"fmt"
	"io"
)

type FileContent struct {
	chunk  ChunkInfo
	buffer io.Reader
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

func (this *Buffer) transformReaders(chunked map[int64]FileContent) []FileContent {
	arr := make([]FileContent, 0)

	done := make(map[int64]bool)

	for _, value := range chunked {
		if _, ok := done[value.chunk.ChunkEnd]; ok {
			continue
		}
		fContent, ok := chunked[value.chunk.ChunkEnd]
		if ok {
			chunkInfo := ChunkInfo{
				ChunkStart: value.chunk.ChunkStart,
				ChunkEnd:   fContent.chunk.ChunkEnd,
				ChunkStr:   "",
			}
			buffer := io.MultiReader(value.buffer, fContent.buffer)
			arr = append(arr, FileContent{chunk: chunkInfo, buffer: buffer})
			done[chunkInfo.ChunkEnd] = true
		} else {
			arr = append(arr, value)
		}
	}

	return arr
}

func (this *Buffer) Flush() {
	chunkedContent := make(map[int64]FileContent)
	for i := range this.buffer {
		chunkedContent[i.chunk.ChunkStart] = i
	}

	joinedBuffers := this.transformReaders(chunkedContent)
	for _, val := range joinedBuffers {
		this.writer.Write(val.chunk.ChunkStart, val.buffer)
	}
}

func (this *Buffer) Close() {
	close(this.buffer)
	if len(this.buffer) > 0 {
		this.Flush()
	}
	this.writer.Close()
}
