# gocrawl
Simple Web Crawler written in Go.

The crawler is limited to one domain - when you start with https://example.com/, it would crawl all pages within example.com, but not follow external links.
Given a URL, it prints a simple site map, showing the links between pages.


## How to use

Testing:
```
make test
```

Formatting:
```
make fmt
```

Linting:
```
make lint
```

Building and running:
```
go build
./gocrawl --help
./gocrawl -n 1000 -d https://test.com/ > sitemap
```
