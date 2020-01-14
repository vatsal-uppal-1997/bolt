package Worker

import (
	"net/http"
	"sync"

	"github.com/vatsal-uppal-1997/bolt/Error"
)

type Worker struct {
	Url    string
	Client *http.Client

	ContentLength int64
	Wg            *sync.WaitGroup

	Dumper *Buffer
}

func (this *Worker) Init(fileLength int64, fileName string, bufferSize int) {
	this.Client = &http.Client{}
	this.Wg = &sync.WaitGroup{}
	writer := &FileWriter{}
	writer.FileLength = fileLength
	writer.FileName = fileName
	this.Dumper = &Buffer{}
	this.Dumper.Init(bufferSize, writer)
}

func (this *Worker) Work(offset int64, sliceRange string) {
	defer this.Wg.Done()
	req, error := http.NewRequest("GET", this.Url, nil)
	Error.HandleError(error, "Error while crafting request", true)
	req.Header.Set("range", "bytes="+sliceRange)

	resp, error := this.Client.Do(req)
	Error.HandleError(error, "Request failed", true)

	this.Dumper.Done(FileContent{offset: offset, buffer: resp.Body})
}
