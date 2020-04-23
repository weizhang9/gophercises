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
	fmt.Println(u.Path)
	body, err := s.getEntryBody()
	if err != nil {
		log.Fatalln(err)
	}
	s.Links = s.getLinks(body)
	if s.contains(s.Entry) {
		s.Links = append(s.Links, s.Entry)
	}
	return s
}

func (s Sitemap) getEntryBody() (string, error) {
	resp, err := http.Get(s.Entry)
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

func (s Sitemap) getLinks(body string) []string {
	node, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatalln(err)
	}
	atags := getATags(node)
	// var links []string
	for _, t := range atags {
		href := getHref(t)
		if s.selfURL(href) && !s.contains(href) {
			s.Links = append(s.Links, href)
		}
	}
	return s.Links
}

func (s Sitemap) selfURL(u string) bool {
	if len(u) > 1 && strings.HasPrefix(u, "/") {
		return true
	}
	if strings.HasPrefix(u, s.Domain) {
		return true
	}
	return false
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
		// base is absolute url
		if strings.HasPrefix(v, s.Domain) {
			// el is absolute url with same domain
			if strings.HasPrefix(el, s.Domain) && v == el {
				return true
			}

			// el is relative url, strip base
			u, _ := url.Parse(v)
			if u.Path == el {
				return true
			}
		}

		// base and el are both relative urls
		if !strings.HasPrefix(el, s.Domain) && v == el {
			return true
		}

		u, err := url.Parse(el)
		if err != nil {
			log.Fatalln(err)
		}
		// can't guarantee selfURL will be run first
		// if it's an external url, throw it away
		if fmt.Sprintf("%s://%s", u.Scheme, u.Host) == s.Domain && u.Path == v {
			return true
		}
	}

	return false
}