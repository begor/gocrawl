package parser

import (
	"io"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)


func ExtractLinksWithCurrentHost(current_url *url.URL, reader io.Reader) []string {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Fatal("Unable to read document")
	}

	links := make([]string, 0)

	for _, node := range doc.Find("a").Nodes {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			href := attr.Val

			href_url, err := url.Parse(href)

			if err != nil {
				log.Fatal("Unable to parse URL ", href)
			}

			full_url := current_url.ResolveReference(href_url)

			if full_url.Host != current_url.Host {
				continue
			}

			full_url.Fragment = ""

			links = append(links, full_url.String())
		}
	}

	return links
}
