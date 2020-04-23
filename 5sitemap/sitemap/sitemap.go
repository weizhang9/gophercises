package sitemap

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Sitemap struct contains everything you need for to build a sitemap
type Sitemap struct {
	Entry string
	Domain string
	Links []string
}

// Parse func takes an entry point url and return a Sitemap struct
func Parse(entryurl string) Sitemap {
	u, err := url.Parse(entryurl)
	if err != nil {
		log.Fatalln(err)
	}
	var s Sitemap
	s.Entry = entryurl
	s.Domain = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	body, err := s.getBody(s.Entry)
	if err != nil {
		log.Fatalln(err)
	}
	s.Links = s.getLinks(body)
	if !s.contains(s.Entry) {
		s.Links = append(s.Links, s.Entry)
	}
	for i := 0; i < 3; i++ {
		s.scapeLinks(s.Links)
	}

	return s
}

func (s *Sitemap) scapeLinks(urls []string) {
	for _, u := range urls {
		body, err := s.getBody(u)
		if err != nil {
			log.Fatalln(err)
		}
		s.Links = s.getLinks(body)
	}
}

func (s Sitemap) getBody(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body),  nil
}

func (s *Sitemap) getLinks(body string) []string {
	node, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatalln(err)
	}
	atags := getATags(node)
	// var links []string
	for _, t := range atags {
		href := getHref(t)
		if s.internalURL(href) && !s.absoluteURL(href) {
			// make sure all urls are absolute path
			href = s.Domain + href
			if !s.contains(href) {
				s.Links = append(s.Links, href)
			}
		}
	}
	return s.Links
}

func (s Sitemap) internalURL(u string) bool {
	if len(u) > 1 && strings.HasPrefix(u, "/") {
		return true
	}
	if s.absoluteURL(u) {
		return true
	}
	return false
}

// absolute internal url
func (s Sitemap) absoluteURL(u string) bool {
	return strings.HasPrefix(u, s.Domain)
}

func getHref(n *html.Node) string {
	var l string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			l = attr.Val
		}
	}
	return l
}

func getATags(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var as []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		as = append(as, getATags(c)...)
	}
	return as
}

func (s Sitemap) contains(el string) bool {
	for _, v := range s.Links {
		if v == el {
			return true
		}
	}

	return false
}