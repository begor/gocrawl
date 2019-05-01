package crawler

import (
	"errors"
	"gocrawl/sitemap"
	"io"
	"strings"
	"testing"
)

type InMemoryFetcher struct {
	pages  map[string]io.Reader
	errors map[string]error
}

func (fetcher *InMemoryFetcher) Get(url string) (io.Reader, error) {
	if val, ok := fetcher.errors[url]; ok {
		return nil, val
	}

	if val, ok := fetcher.pages[url]; ok {
		return val, nil
	}

	return nil, errors.New("Not found")
}

func (fetcher *InMemoryFetcher) AddUrlToPages(url string, reader io.Reader) {
	fetcher.pages[url] = reader
}

func (fetcher *InMemoryFetcher) AddUrlToErrors(url string, err error) {
	fetcher.errors[url] = err
}

func createInMemoryFetcher() InMemoryFetcher {
	return InMemoryFetcher{pages: make(map[string]io.Reader), errors: make(map[string]error)}
}

// E2E-like test with InMemoryFetcher.
// Having DI for Fetcher in Crawler, we can easilly test functionality without mocks.
// More specific tests can be found in parser/sitemap packages.
func TestEndToEndCrawlingWithMockFetcher(t *testing.T) {
	sitemap := sitemap.Create()
	fetcher := createInMemoryFetcher()
	// "Root" page with two links to /page/ and /another-page/
	fetcher.AddUrlToPages("https://test.com/", strings.NewReader("<a href=\"/page/\"></a><a href=\"/another-page/\"></a>"))
	// Level-1 /page/  with one new page
	fetcher.AddUrlToPages("https://test.com/page/", strings.NewReader("<a href=\"new-page/\"></a>"))
	// Level-2 /page/new-page/ page with one link to foreign host and one already visited page (because it's child of root)
	fetcher.AddUrlToPages("https://test.com/page/new-page/", strings.NewReader("<a href=\"https://not-test.com/\"></a><a href=\"/another-page/\"></a>"))
	// Level-3 /another-page/ child page with link to error page and one visited root
	fetcher.AddUrlToPages("https://test.com/another-page/", strings.NewReader("<a href=\"42/\"></a><a href=\"..\"></a>"))
	fetcher.AddUrlToErrors("https://test.com/another-page/42/", errors.New("42"))

	Crawl("https://test.com/", &fetcher, &sitemap)

	if sitemap.Size() != 4 {
		t.Error("sitemap.Size(): expected 4, got", sitemap.Size())
	}

	if sitemap.ErrSize() != 1 {
		t.Error("sitemap.ErrSize()", "expected 1, got", sitemap.ErrSize())
	}

	if links, ok := sitemap.GetLinks("https://test.com/"); !ok || len(links) != 2 {
		t.Error("https://test.com/", "expected 2, got", len(links))
	}

	if links, ok := sitemap.GetLinks("https://test.com/page/"); !ok || len(links) != 1 {
		t.Error("https://test.com/page/", "expected 1, got", len(links))
	}

	// NOTE: no link to foreign host
	if links, ok := sitemap.GetLinks("https://test.com/page/new-page/"); !ok || len(links) != 1 {
		t.Error("https://test.com/page/new-page/", "expected 1, got", len(links))
	}

	// NOTE: including link to error page
	if links, ok := sitemap.GetLinks("https://test.com/another-page/"); !ok || len(links) != 2 {
		t.Error("https://test.com/another-page/", "expected 2, got", len(links))
	}

	// NOTE: foreign host
	if sitemap.Contains("https://not-test.com/") {
		t.Error("https://not-test.com/", "expected false, got true")
	}

	// NOTE: error page is not present
	if sitemap.Contains("https://test.com/another-page/42/") {
		t.Error("https://test.com/another-page/42/", "expected false, got true")
	}
}
