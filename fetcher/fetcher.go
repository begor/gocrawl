package fetcher

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("fetcher")

type Fetcher interface {
	Get(string) (io.Reader, error)
}

type HTTPLimitFetcher struct {
	sem_chan chan int
}

func CreateHTTPFetcher(limit int) HTTPLimitFetcher {
	sem_chan := make(chan int, limit)

	return HTTPLimitFetcher{sem_chan: sem_chan}
}

func (fetcher *HTTPLimitFetcher) Get(url string) (io.Reader, error) {
	fetcher.sem_chan <- 1

	res, err := http.Get(url)

	<-fetcher.sem_chan

	if err != nil {
		log.Warning("Unable to fetch", url, err)
		return nil, err
	}

	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)

	return bytes.NewReader(bs), err
}
