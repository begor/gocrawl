package crawler

import (
	"gocrawl/fetcher"
	"gocrawl/parser"
	"gocrawl/sitemap"
	"net/url"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("crawler")


// Main crawler routine: 
// - keeps visited set in sitemap
// - fetches byte stream via given fetcher.Fetcher (allows us to use DI, see crawler_test.go for example)
// - extracts outgoing links from byte stream with parser.ExtractLinksWithCurrentHost
// - traverses (in DFS-way) outgoing links from byte stream and runs Crawl for them concurrently
func Crawl(url_str string, fetcher fetcher.Fetcher, sitemap *sitemap.Sitemap) {
	log.Debug("Crawling: ", url_str)

	if sitemap.Contains(url_str) {
		log.Debug(url_str, " already visited")
		return
	}

	sitemap.MarkVisited(url_str)

	current_url, err := url.Parse(url_str)

	if err != nil {
		log.Warning("Unable to get URL of ", url_str, err)

		sitemap.SetError(url_str, err)
		return
	}

	response, err := fetcher.Get(url_str)

	if err != nil {
		log.Warning("Unable to fetch from ", url_str, err)

		sitemap.SetError(url_str, err)
		return
	}

	links, err := parser.ExtractLinksWithCurrentHost(current_url, response)

	if err != nil {
		log.Warning("Unable to extract links from ", url_str, err)

		sitemap.SetError(url_str, err)
		return
	}

	sitemap.SetLinks(url_str, links)

	crawlChilderLinks(links, fetcher, sitemap)
}

func crawlChilderLinks(links []string, fetcher fetcher.Fetcher, sitemap *sitemap.Sitemap) {
	children_queue := make(chan bool)

	for _, link := range links {
		go func(link string) {
			Crawl(link, fetcher, sitemap)
			children_queue <- true
		}(link)
	}

	for range links {
		<-children_queue
	}
}
