package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	entry := flag.String("url", "", "main entrypoint for sitemap builder")
	flag.Parse()
	if !strings.HasPrefix(*entry, "http") {
		log.Fatalln("Please provide an entrypoint url with protocol scheme")
	}
	body, err := getEntryBody(*entry)
	if err != nil {
		log.Fatalln(err)
	}
	links := getLinks(body, *entry)
	fmt.Println(links)
}

func getEntryBody(url string) (string, error) {
	resp, err := http.Get(url)
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

func getLinks(body, entryurl string) []string {
	node, err := html.Parse(strings.NewReader(body))
	if err != nil {
		log.Fatalln(err)
	}
	atags := getATags(node)
	var links []string
	for _, t := range atags {
		href := getHref(t)
		links = append(links, validateUrls(entryurl, uniqueUrl(links, href)))
	}
	return links
}

func uniqueUrl(l []string, url string) string {
	if contains(l, url) {
		return ""
	}
	return url
}

func validateUrls(domain, uniqueUrl string) string {
	if len(uniqueUrl) > 1 && strings.HasPrefix(uniqueUrl, "/") {
		return uniqueUrl
	}
	if strings.HasPrefix(uniqueUrl, domain) {
		return uniqueUrl
	}
	return ""
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

func contains(base []string, el string) bool {
	for _, v := range base {
		if v == el {
			return true
		}
	}

	return false
}