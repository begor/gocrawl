package parser

import (
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func ExtractLinksWithCurrentHost(current_url *url.URL, reader io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		// TODO: log.Warn
		return nil, err
	}

	// Keep links unique to reduce extra work on fetching later
	links_map := make(map[string]bool)

	for _, node := range doc.Find("a").Nodes {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			href := attr.Val

			href_url, err := url.Parse(href)

			if err != nil {
				// TODO: log.Warn
				continue
			}

			full_url := current_url.ResolveReference(href_url)

			if full_url.Host != current_url.Host {
				continue
			}

			full_url.Fragment = ""
			full_url.RawQuery = ""

			full_url_s := full_url.String()

			if _, ok := links_map[full_url_s]; ok {
				continue
			}

			links_map[full_url_s] = true
		}
	}

	links := make([]string, 0)

	for link, _ := range links_map {
		links = append(links, link)
	}

	return links, nil
}
