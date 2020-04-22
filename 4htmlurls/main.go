package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	htmls, err := ioutil.ReadDir("htmls")
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range htmls {
		data, err := ioutil.ReadFile(fmt.Sprintf("htmls/%s",h.Name()))

		if err != nil {
			log.Fatal(err)
		}

		nodes, err := html.Parse(bytes.NewReader(data))
		if err != nil {
			log.Fatal(err)
		}
		links := getLinks(nodes)
		fmt.Printf("%+v", links)
	}
}

func getLinks(nodes *html.Node) []Link {
	atags := getATags(nodes)

	var links []Link
	for _, a := range atags {
		links = append(links, buildLink(a))
	}

	return links
}

func getATags(node *html.Node) []*html.Node {
	var aTags []*html.Node
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		aTags = append(aTags, getATags(c)...)
	}

	return aTags
}

func buildLink(n *html.Node) Link {
	var l Link
	// get href
	for _, v := range n.Attr {
		if v.Key == "href" {
			l.Href = v.Val
		}
		break
	}

	// get text
	l.Text = getText(n)

	return l
}

func getText(anode *html.Node) string {
	var t string
	if anode.Type == html.TextNode {
		return anode.Data
	}

	for c := anode.FirstChild; c != nil; c = c.NextSibling {
		t = strings.Join(strings.Fields(t+getText(c)), " ")
	}

	return t
}
