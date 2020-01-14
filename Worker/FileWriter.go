package Worker

import (
	"fmt"
	errorHandler "github.com/vatsal-uppal-1997/bolt/Error"
	"io"
	"os"
	"sync"
)

type FileWriter struct {
	Mutex      sync.Mutex
	Downloaded float64
	FileLength int64
	FileName   string
	Fd         *os.File
}

func (this *FileWriter) Close() {
	this.Fd.Close()
}

func (this *FileWriter) Write(offset int64, content io.ReadCloser) {
	_, err := os.Stat(this.FileName)
	this.Mutex.Lock()

	if os.IsNotExist(err) {
		this.Fd, err = os.Create(this.FileName)
		errorHandler.HandleError(err, "Unable to create new file", true)
		this.Fd.Truncate(this.FileLength + 100)
	}

	_, err = this.Fd.Seek(offset, 0)
	errorHandler.HandleError(err, "Unable to seek to file", true)

	bytes, err := io.Copy(this.Fd, content)
	this.Fd.Seek(0, 0)
	errorHandler.HandleError(err, "Unable to copy buffer", true)

	this.Downloaded += float64(bytes)
	donePercentage := this.Downloaded / float64(this.FileLength) * 100
	fmt.Printf("%f percent file downloaded\n", donePercentage)
	this.Mutex.Unlock()
}

func (this *FileWriter) WriteSync(content io.ReadCloser) {
	_, err := os.Stat(this.FileName)

	this.Fd, err = os.Create(this.FileName)
	errorHandler.HandleError(err, "Unable to create new file", true)
	io.Copy(this.Fd, content)
}
