package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/vatsal-uppal-1997/bolt/Error"
	"github.com/vatsal-uppal-1997/bolt/Worker"
)

func supportsThreading(url string) (bool, int64) {
	resp, err := http.Head(url)
	Error.HandleError(err, "Head request not supported", true)
	contentLengthStr, acceptRanges := resp.Header.Get("content-length"), resp.Header.Get("Accept-Ranges")

	contentLength, err := strconv.Atoi(contentLengthStr)
	Error.HandleError(err, "Cannot convert contentLength to integer", true)

	if acceptRanges != "bytes" {
		return false, 0 // synchronus download
	}

	return true, int64(contentLength)
}

func threaded(url string, contentLength int64, workers int, fileName string) {
	chunks := int(contentLength) / workers

	var worker Worker.Worker

	worker.Url = url
	worker.Init(int64(contentLength), fileName)
	for i := 0; i < int(contentLength); i += chunks {
		worker.Wg.Add(1)
		if i+chunks >= int(contentLength) {
			go worker.Work(int64(i), strconv.Itoa(i)+"-"+strconv.Itoa(int(contentLength)))
			break
		}
		go worker.Work(int64(i), strconv.Itoa(i)+"-"+strconv.Itoa(i+chunks))
	}
	worker.Wg.Wait()
	worker.Dumper.Close()
}

func synchronus(url, fileName string) {
	resp, err := http.Get(url)

	Error.HandleError(err, "Some error occured while downloading the file", true)
	writer := Worker.FileWriter{}
	writer.FileName = fileName
	writer.WriteSync(resp.Body)
}

func main() {

	workers := flag.Int("workers", 10, "Number of workers")
	url := flag.String("url", "http://www.blabla.com", "The url of file to download")
	out := flag.String("out", "bolt.jpeg", "Name of the output file")

	flag.Parse()

	threading, contentLength := supportsThreading(*url)

	if threading {
		threaded(*url, contentLength, *workers, *out)
	} else {
		log.Println("Multi-Threaded download not supported - downloading synchronously")
		synchronus(*url, *out)
	}

}
