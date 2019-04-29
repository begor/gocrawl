package parser

import (
	"net/url"
	"strings"
	"testing"
)

func TestReaderWithoutLinksString(t *testing.T) {
	reader := strings.NewReader("This is a page without links")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)

	if len(links) != 0 {
		t.Error("Expected 0, got ", len(links))
	}
}

func TestReaderWithoutLinksHTML(t *testing.T) {
	reader := strings.NewReader("<ab></ab>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)

	if len(links) != 0 {
		t.Error("Expected 0, got ", len(links))
	}
}

func TestReaderOneLinkDifferentHost(t *testing.T) {
	reader := strings.NewReader("<a href=\"https://not-a-test.com/page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)

	if len(links) != 0 {
		t.Error("Expected 0, got ", len(links))
	}
}

func TestReaderOneLinkSameHost(t *testing.T) {
	reader := strings.NewReader("<a href=\"https://test.com/page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link := "https://test.com/page/"

	if len(links) != 1 || links[0] != expected_link {
		t.Error("Expected ", expected_link, " got ", links)
	}
}

func TestReaderOneLinkRelativeSimple(t *testing.T) {
	reader := strings.NewReader("<a href=\"/page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link := "https://test.com/page/"

	if len(links) != 1 || links[0] != expected_link {
		t.Error("Expected ", expected_link, " got ", links)
	}
}

func TestReaderOneLinkRelativeComplex(t *testing.T) {
	reader := strings.NewReader("<a href=\"/page/temp/../42/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link := "https://test.com/page/42/"

	if len(links) != 1 || links[0] != expected_link {
		t.Error("Expected ", expected_link, " got ", links)
	}
}

func TestReaderTwoLinksWithSameHost(t *testing.T) {
	reader := strings.NewReader("<a href=\"https://test.com/page/\"></a><a href=\"https://test.com/another/page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link1 := "https://test.com/page/"
	expected_link2 := "https://test.com/another/page/"

	if len(links) != 2 || links[0] != expected_link1 || links[1] != expected_link2 {
		t.Error("Expected ", expected_link1, expected_link2, " got ", links)
	}
}

func TestReaderTwoLinksRelative(t *testing.T) {
	reader := strings.NewReader("<a href=\"/page/\"></a><a href=\"/another/42/../page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link1 := "https://test.com/page/"
	expected_link2 := "https://test.com/another/page/"

	if len(links) != 2 || links[0] != expected_link1 || links[1] != expected_link2 {
		t.Error("Expected ", expected_link1, expected_link2, " got ", links)
	}
}

func TestReaderTwoLinksOneDifferentHost(t *testing.T) {
	reader := strings.NewReader("<a href=\"https://not-a-test.com/page/\"></a><a href=\"/another/42/../page/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)
	expected_link := "https://test.com/another/page/"

	if len(links) != 1 || links[0] != expected_link {
		t.Error("Expected ", expected_link, " got ", links)
	}
}

func TestReaderTwoLinksTwoDifferentHost(t *testing.T) {
	reader := strings.NewReader("<a href=\"https://not-a-test.com/page/\"></a><a href=\"https://not-a-test.com/\"></a>")
	url := url.URL{Host: "test.com", Scheme: "https"}

	links := ExtractLinksWithCurrentHost(&url, reader)

	if len(links) != 0 {
		t.Error("Expected 0 got ", links)
	}
}
