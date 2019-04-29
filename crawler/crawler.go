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
// 0) if given url_str already in sitemap, return
// 1) fetch document from url_str if fetch_limiter is available or block
// 2) parse links from document using parser package
// 3) run DFS on links using crawlChilderLinks
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

		sitemap.AddError(url_str, err)
		return
	}

	response, err := fetcher.Get(url_str)

	if err != nil {
		log.Warning("Unable to fetch from ", url_str, err)

		sitemap.AddError(url_str, err)
		return
	}

	links, err := parser.ExtractLinksWithCurrentHost(current_url, response)

	if err != nil {
		log.Warning("Unable to extract links from ", url_str, err)

		sitemap.AddError(url_str, err)
		return
	}

	sitemap.AddLinks(url_str, links)

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
