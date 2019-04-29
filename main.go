package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"gocrawl/crawler"
	"gocrawl/fetcher"
	"gocrawl/sitemap"
	"time"
)

var log = logging.MustGetLogger("main")

func setupLogging(debug bool) {
	var format = logging.MustStringFormatter(
		`%{color}%{time} %{shortfunc} ▶ %{level:.1s} %{color:reset} %{message}`,
	)

	logging.SetFormatter(format)

	if debug {
		logging.SetLevel(logging.DEBUG, "")
	} else {
		logging.SetLevel(logging.INFO, "")
	}
}

func main() {
	http_timeout := flag.Int("t", 10, "Timeout in seconds for HTTP fetching")
	parallel_fetchers := flag.Int("n", 32, "Number of parallel coroutines")
	short_view := flag.Bool("s", false, "Short view (without page links)")
	debug := flag.Bool("d", false, "Enable verbose DEBUG logging")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("gocrawl is a simple Web Crawler.")
		fmt.Println("Usage:\n\tgocrawl <options> [site url]")
		return
	}

	url := flag.Args()[0]

	setupLogging(*debug)

	sitemap := sitemap.Create()
	fetcher := fetcher.CreateHTTPFetcher(*parallel_fetchers, *http_timeout)

	log.Info("Starting Crawler with", *parallel_fetchers, "concurrent fetchers")

	start := time.Now()
	crawler.Crawl(url, &fetcher, &sitemap)
	fin := time.Now()

	elapsed := fin.Sub(start)

	log.Info("Crawler done at", elapsed, "seconds!")

	sitemap.PrintReport(*short_view)
}
