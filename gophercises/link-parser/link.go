package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	docNode, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := dfs(docNode)
	links := make([]Link, 0)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) Link {
	link := Link{}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}
	link.Text = text(n)
	return link
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(strings.Trim(n.Data, "\n"))
	}
	if n.Type != html.ElementNode {
		return ""
	}
	res := ""
	for curr := n.FirstChild; curr != nil; curr = curr.NextSibling {
		res += text(curr) + " "
	}
	return strings.Trim(res, " ")
}

func dfs(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var res []*html.Node
	for curr := n.FirstChild; curr != nil; curr = curr.NextSibling {
		res = append(res, dfs(curr)...)
	}
	return res
}
