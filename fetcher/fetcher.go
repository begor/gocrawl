package fetcher

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("fetcher")

// Using interface allows us to use DI and change implementation from HTTP to in-memory/file/...
type Fetcher interface {
	Get(string) (io.Reader, error)
}

// NOTE: in "real production" that should include retry/backoff policies, metering and so on.
type HTTPLimitFetcher struct {
	sem_chan chan int
	client   *http.Client
}

func CreateHTTPFetcher(limit int, timeout_sec int) HTTPLimitFetcher {
	sem_chan := make(chan int, limit)
	timeout := time.Duration(time.Duration(timeout_sec) * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	return HTTPLimitFetcher{sem_chan: sem_chan, client: &client}
}

func (fetcher *HTTPLimitFetcher) Get(url string) (io.Reader, error) {
	fetcher.sem_chan <- 1

	res, err := fetcher.client.Get(url)

	<-fetcher.sem_chan

	if err != nil {
		log.Warning("Unable to fetch", url, err)
		return nil, err
	}

	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)

	return bytes.NewReader(bs), err
}
