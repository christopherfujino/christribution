package main

import (
	"fmt"
	"strings"

	"net/http"

	html "golang.org/x/net/html"
)

const manifest = "https://www.linuxfromscratch.org/lfs/view/development/chapter03/packages.html"

func main() {
	res, err := http.Get(manifest)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		panic(fmt.Sprintf("Request for %s failed with code %d", manifest, res.StatusCode))
	}
	rootNode, err := html.Parse(res.Body)
	if err != nil {
		panic("Failed to parse HTML content")
	}
	archives := traverse(rootNode)
	for _, archive := range archives {
		fmt.Printf("Found: %s\n", archive)
	}
}

func traverse(node *html.Node) []string {
	var archives []string

	for node := range node.ChildNodes() {
		if node.Type == html.ElementNode && node.Data == "dl" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "variablelist" {
					archives = append(archives, findChildParagraph(node)...)
					break
				}
			}
		} else {
			archives = append(archives, traverse(node)...)
		}
	}
	return archives
}

func findChildParagraph(node *html.Node) []string {
	var urls []string

	if node.Type == html.ElementNode && node.Data == "p" {
		var string_opt *string
		for child := range node.ChildNodes() {
			if child.Type == html.TextNode {
				var data = strings.TrimSpace(child.Data)
				if data == "Download:" {
					string_opt = findChildAnchor(node)
					if string_opt == nil {
						panic("Oops")
					}
					fmt.Printf("Found %s\n", *string_opt)
					urls = append(urls, *string_opt)
				} else {
					fmt.Printf("Got text node: \"%s\"\n", data)
				}
			}
		}
	} else {
		for child := range node.ChildNodes() {
			urls = append(urls, findChildParagraph(child)...)
		}
	}
	return urls
}

func findChildAnchor(node *html.Node) *string {
	fmt.Printf("Looking for an anchor in %s\n", node.Data)
	if node.Type == html.ElementNode && node.Data == "a" {
		fmt.Println("Find an anchor...")
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return &attr.Val
			}
		}
		fmt.Println("...but it didn't have an href attribute.")
	} else {
		var string_opt *string
		for child := range node.ChildNodes() {
			string_opt = findChildAnchor(child)
			if string_opt != nil {
				return string_opt
			}
		}
	}

	return nil
}
