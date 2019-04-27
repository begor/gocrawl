package main

import (
	"fmt"
	"log"
	"path"

	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	host := "https://google.com"
	resp, err := http.Get(host)
	
	if err != nil {
		log.Fatal("Unable to fetch")
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal("Unable to parse")
	}

	urls := doc.Find("a").Map(func(i int, s *goquery.Selection) (string) {
		url, _ := s.Attr("href")

		return url
	})

	for _, url_str := range urls {
		fmt.Printf("%s ", url_str)
		u, _ := url.Parse(url_str)
		
		if !(u.IsAbs()) {
			host_u, _ := url.Parse(host)
			host_u.Path = path.Join(host_u.Path, u.Path)
			u = host_u
		}

		fmt.Println(u.String())
	}
}