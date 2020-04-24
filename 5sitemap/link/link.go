package link

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link struct
type Link struct {
	Href string
	Text string
}

// Parse main API func
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	// tree visual
	// dfs(doc, "")

	// find <a> nodes
	nodes := linkNodes(doc)
	var links []Link
	// for each <a> node, build a Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	// return []Link
	return links, nil
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}

	ret.Text = getText(n)
	return ret
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func dfs(n *html.Node, p string) {
	msg := n.Data
	if n.Type == html.ElementNode {
		msg = "<" + msg + ">"
	}
	if n.Type == html.TextNode {
		msg = "*\"" + msg + "\"*"
	}
	if n.Type == html.DocumentNode {
		msg = "^BOSS" + msg + "^"
	}
	if n.Type == html.CommentNode {
		msg = "###" + msg + "###"
	}
	if n.Type == html.DoctypeNode {
		msg = "??" + msg + "??"
	}
	fmt.Println(p, msg)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dfs(c, p+"£")
	}
}
