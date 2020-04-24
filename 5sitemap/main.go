package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"sitemap/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

func main() {
	entry := flag.String("url", "https://gophercises.com", "entry url to build the sitemap for")
	depth := flag.Int("depth", 3, "max traverse depth")
	flag.Parse()
	urls := bfs(*entry, *depth)
	
	// encode urls into xml sitemap
	sitemap := urlset {
		Xmlns: xmlns,
	}
	sitemap.Url = make([]loc, len(urls))
	for i, u := range urls {
		sitemap.Url[i] = loc{u}
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	fmt.Print(xml.Header)
	enc.Encode(sitemap)
}

type urlset struct {
	Url []loc `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

type loc struct {
	Loc string `xml:"loc"`
}

func bfs(entry string, depth int) []string {
	seen := make(map[string]struct{})
	q := make(map[string]struct{})
	nq := map[string]struct{}{
		entry: struct{}{},
	}

	for i := 0; i <= depth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, nqUrl := range scrape(url) {
				nq[nqUrl] = struct{}{}
			}
		}
	}

	urls := make([]string, 0, len(seen))
	for u := range seen {
		urls = append(urls, u)
	}
	return urls
}

func scrape(entry string) []string {
	resp, err := http.Get(entry)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return getUrls(resp.Body, base)
}

func getUrls(html io.Reader, base string) []string {
	var urls []string
	links, _ := link.Parse(html)
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			href := base + l.Href
			if !contains(urls, href) {
				urls = append(urls, href)
			}
		case strings.HasPrefix(l.Href, base):
			if !contains(urls, l.Href) {
				urls = append(urls, l.Href)
			}
		}
	}

	return urls
}

func contains(urls []string, u string) bool {
	for _, v := range urls {
		if v == u {
			return true
		}
	}
	return false
}
