package fetcher

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Fetcher interface {
	Get(string) io.Reader
}


type HTTPLimitFetcher struct {
	sem_chan chan int
}


func CreateHTTPFetcher(limit int) HTTPLimitFetcher {
	sem_chan := make(chan int, limit)

	return HTTPLimitFetcher{sem_chan: sem_chan}
}


func (fetcher *HTTPLimitFetcher) Get(url string) io.Reader {
	fetcher.sem_chan <- 1

	res, err := http.Get(url)

	<-fetcher.sem_chan

	if err != nil {
		log.Fatal("Unable to fetch ", url)
	}

	defer res.Body.Close()

	bs, _ := ioutil.ReadAll(res.Body)

	return bytes.NewReader(bs)
}