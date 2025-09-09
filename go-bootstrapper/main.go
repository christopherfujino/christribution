package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"net/http"

	html "golang.org/x/net/html"
)

const manifest = "https://www.linuxfromscratch.org/lfs/view/development/chapter03/packages.html"

func main() {
	flag.Parse()
	var args = flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "bootstrap":
			bootstrap()
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown sub-command: %s\n", args[0])
			flag.Usage()
			os.Exit(1)
		}
	} else {
		flag.Usage()
		os.Exit(0)
	}
}

func bootstrap() {
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

	archiveLength := len(archives)
	var archive string
	fmt.Println("[")
	for i := 0; i < archiveLength; i += 1 {
		archive = archives[i]
		if i == archiveLength - 1 {
			fmt.Printf("  \"%s\"\n", archive)
		} else {
			fmt.Printf("  \"%s\",\n", archive)
		}
	}
	fmt.Println("]")
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
					urls = append(urls, *string_opt)
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
	fmt.Fprintf(os.Stderr, "Looking for an anchor in %s\n", node.Data)
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return &attr.Val
			}
		}
		fmt.Fprintf(os.Stderr, "found anchor but it didn't have an href attribute.")
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
