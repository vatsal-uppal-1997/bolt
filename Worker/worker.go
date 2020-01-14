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

	Dumper *FileWriter
}

func (this *Worker) Init(fileLength int64, fileName string) {
	this.Client = &http.Client{}
	this.Wg = &sync.WaitGroup{}
	this.Dumper = &FileWriter{}
	this.Dumper.FileLength = fileLength
	this.Dumper.FileName = fileName
}

func (this *Worker) Work(offset int64, sliceRange string) {
	defer this.Wg.Done()
	req, error := http.NewRequest("GET", this.Url, nil)
	Error.HandleError(error, "Error while crafting request", true)
	req.Header.Set("range", "bytes="+sliceRange)

	resp, error := this.Client.Do(req)
	Error.HandleError(error, "Request failed", true)

	defer resp.Body.Close()
	this.Dumper.Write(offset, resp.Body)
}
