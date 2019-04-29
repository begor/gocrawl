package sitemap

import (
	"errors"
	"testing"
)

func TestSizeEmpty(t *testing.T) {
	sitemap := Create()

	if sitemap.Size() != 0 {
		t.Error("Expected 0, got", sitemap.Size())
	}
}

func TestErrSizeEmpty(t *testing.T) {
	sitemap := Create()

	if sitemap.ErrSize() != 0 {
		t.Error("Expected 0, got", sitemap.ErrSize())
	}
}

func TestMarkVisitedOneSize(t *testing.T) {
	sitemap := Create()
	sitemap.MarkVisited("test")

	if sitemap.Size() != 1 {
		t.Error("Expected 1, got", sitemap.Size())
	}
}

func TestMarkVisitedOneErrSize(t *testing.T) {
	sitemap := Create()
	sitemap.MarkVisited("test")

	if sitemap.ErrSize() != 0 {
		t.Error("Expected 0, got", sitemap.ErrSize())
	}
}

func TestMarkVisitedOneContains(t *testing.T) {
	sitemap := Create()
	sitemap.MarkVisited("test")

	if !sitemap.Contains("test") {
		t.Error("Expected true, got", sitemap.ErrSize())
	}
}

func TestMarkVisitedTwoSameSize(t *testing.T) {
	sitemap := Create()
	sitemap.MarkVisited("test")
	sitemap.MarkVisited("test")

	if sitemap.Size() != 1 {
		t.Error("Expected 1, got", sitemap.Size())
	}
}

func TestMarkVisitedTwoDistinctSize(t *testing.T) {
	sitemap := Create()
	sitemap.MarkVisited("test")
	sitemap.MarkVisited("test2")

	if sitemap.Size() != 2 {
		t.Error("Expected 2, got", sitemap.Size())
	}
}

func TestSetLinksOne(t *testing.T) {
	sitemap := Create()
	sitemap.SetLinks("test", []string{"test2"})

	links, _ := sitemap.GetLinks("test")

	if sitemap.Size() != 1 || len(links) != 1 || links[0] != "test2" {
		t.Error("Expected [test2], got", links)
	}
}

func TestSetLinksTwo(t *testing.T) {
	sitemap := Create()
	sitemap.SetLinks("test", []string{"test2", "test3"})

	links, _ := sitemap.GetLinks("test")

	if sitemap.Size() != 1 || len(links) != 2 || links[0] != "test2" || links[1] != "test3" {
		t.Error("Expected [test2, test3], got", links)
	}
}

func TestSetLinksTwoSeparate(t *testing.T) {
	sitemap := Create()
	sitemap.SetLinks("test", []string{"test2"})
	sitemap.SetLinks("test", []string{"test3"})

	links, _ := sitemap.GetLinks("test")

	if sitemap.Size() != 1 || len(links) != 1 || links[0] != "test3" {
		t.Error("Expected [test3], got", links)
	}
}

func TestSetLinksComplex(t *testing.T) {
	sitemap := Create()
	sitemap.SetLinks("test", []string{"test2", "test3"})
	sitemap.SetLinks("test2", []string{"test", "test3"})
	sitemap.SetLinks("test2", []string{"test4"})

	links1, _ := sitemap.GetLinks("test")

	if len(links1) != 2 || links1[0] != "test2" || links1[1] != "test3" {
		t.Error("Expected [test2, test3], got", links1)
	}

	links2, _ := sitemap.GetLinks("test2")

	if len(links2) != 1 || links2[0] != "test4" {
		t.Error("Expected [test4], got", links2)
	}
}

func TestSetErrorFirst(t *testing.T) {
	sitemap := Create()
	sitemap.SetError("test", errors.New("test error"))

	if sitemap.ErrSize() != 1 {
		t.Error("Expected 1, got", sitemap.ErrSize())
	}
}

func TestSetErrorDeleteLinks(t *testing.T) {
	sitemap := Create()
	sitemap.SetLinks("test", []string{"test2", "test3"})
	sitemap.SetError("test", errors.New("test error"))

	if sitemap.ErrSize() != 1 || sitemap.Size() != 0 {
		t.Error("Expected 1, 0, got", sitemap.ErrSize(), sitemap.Size())
	}
}
