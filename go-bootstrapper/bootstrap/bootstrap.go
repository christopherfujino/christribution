package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"net/http"

	html "golang.org/x/net/html"

	"github.com/christopherfujino/christribution/go-bootstrapper/common"
)

const manifest = "https://www.linuxfromscratch.org/lfs/view/stable/chapter03/packages.html"

const patchManifest = "https://www.linuxfromscratch.org/lfs/view/stable/chapter03/patches.html"

func Bootstrap() {
	var archives = findArchives(fetchRemoteNode(manifest))

	archives = rewriteMirrors(archives)

	patches := findArchives(fetchRemoteNode(patchManifest))

	patches = rewriteMirrors(patches)

	jsonBytes, err := json.Marshal(common.Manifest{
		Date:     time.Now(),
		Archives: archives,
		Patches:  patches,
	})
	if err != nil {
		panic(err)
	}
	var indentedBytes = bytes.Buffer{}
	json.Indent(&indentedBytes, jsonBytes, "", "  ")

	outFile, err := os.Create(common.ManifestPath)
	if err != nil {
		panic(err)
	}
	_, err = outFile.Write(indentedBytes.Bytes())

	if err != nil {
		panic(err)
	}
}

type archiveParseState int

const (
	parsingName archiveParseState = iota
	parsingRemoteAndHash
)

func findArchives(node *html.Node) []common.Archive {
	var archives []common.Archive

	for node := range node.ChildNodes() {
		if node.Type == html.ElementNode && node.Data == "dl" {
			for _, attr := range node.Attr {
				if attr.Key == "class" && attr.Val == "variablelist" {
					var state archiveParseState = parsingName
					var name string
					var remote string
					for node := range node.ChildNodes() {
						switch state {
						case parsingName:
							var name_opt = findArchiveName(node.NextSibling)
							if name_opt != nil {
								name = *name_opt
								state = parsingRemoteAndHash
							}
						case parsingRemoteAndHash:
							var remote_opt = findTextFromLabel(
								node,
								"Download:",
								findTextFromAnchor,
							)
							if remote_opt != nil {
								remote = *remote_opt

								var hash_opt = findTextFromLabel(
									node,
									"MD5 sum:",
									findTextFromCode,
								)
								if hash_opt == nil {
									panic("Unreachable")
								}
								archives = append(
									archives,
									common.CreateArchive(name, remote, *hash_opt),
								)
								name = "unreachable"
								remote = "unreachable"
								state = parsingName
							}
						default:
							panic("Unreachable")
						}
					}
					if state != parsingName {
						panic(fmt.Sprintf("Bad state: did not finish parsing %s", name))
					}
				}
			}
		} else {
			archives = append(archives, findArchives(node)...)
		}
	}
	return archives
}

func findArchiveName(node *html.Node) *string {
	for ; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode && node.Data == "dt" {
			for node := range node.ChildNodes() {
				if node.Type == html.ElementNode && node.Data == "span" {
					for _, attr := range node.Attr {
						if attr.Key == "class" && attr.Val == "term" {
							for node := range node.ChildNodes() {
								if node.Type == html.TextNode {
									var returnValue = strings.TrimSpace(node.Data)
									returnValue = strings.TrimSuffix(returnValue, " -")
									return &returnValue
								}
							}
						}
					}
					panic("Unreachable")
				}
			}
		}
	}
	return nil
}

func findTextFromLabel(node *html.Node, label string, predicate func(*html.Node) *string) *string {
	if node.Type == html.ElementNode && node.Data == "p" {
		var string_opt *string
		for child := range node.ChildNodes() {
			if child.Type == html.TextNode {
				var data = strings.TrimSpace(child.Data)
				if data == label {
					string_opt = predicate(node)
					if string_opt == nil {
						panic("Oops")
					}
					return string_opt
				}
			}
		}
	} else {
		for child := range node.ChildNodes() {
			var string_opt = findTextFromLabel(child, label, predicate)
			if string_opt != nil {
				return string_opt
			}
		}
	}
	return nil
}

func findTextFromAnchor(node *html.Node) *string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return &attr.Val
			}
		}
	} else {
		var string_opt *string
		for child := range node.ChildNodes() {
			string_opt = findTextFromAnchor(child)
			if string_opt != nil {
				return string_opt
			}
		}
	}

	return nil
}

func findTextFromCode(node *html.Node) *string {
	if node.Type == html.ElementNode && node.Data == "code" {
		for node := range node.ChildNodes() {
			if node.Type == html.TextNode {
				var returnValue = strings.TrimSpace(node.Data)
				return &returnValue
			}
		}
		panic("Unreachable")
	} else {
		var string_opt *string
		for child := range node.ChildNodes() {
			string_opt = findTextFromCode(child)
			if string_opt != nil {
				return string_opt
			}
		}
	}

	return nil
}

func rewriteMirrors(archives []common.Archive) (outputArchives []common.Archive) {
	const ftpGnuPrefix = "https://ftp.gnu.org/gnu"
	for _, archive := range archives {
		if strings.HasPrefix(archive.Remote, ftpGnuPrefix) {
			archive.Remote = strings.Replace(archive.Remote, ftpGnuPrefix, "https://mirrors.kernel.org/gnu", 1)
		}
		outputArchives = append(outputArchives, archive)
	}

	return
}

func fetchRemoteNode(path string) *html.Node {
	res, err := http.Get(path)
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
	return rootNode
}
