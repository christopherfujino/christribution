package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/christopherfujino/christribution/go-bootstrapper/bootstrap"
	"github.com/christopherfujino/christribution/go-bootstrapper/extract"
	"github.com/christopherfujino/christribution/go-bootstrapper/fetch"
)

func main() {
	flag.Parse()
	var args = flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "bootstrap":
			bootstrap.Bootstrap()
		case "fetch":
			fetch.Fetch()
		case "extract":
			extract.Extract()
		default:
			fmt.Fprintf(os.Stderr, "Unknown sub-command: %s\n", args[0])
			flag.Usage()
			os.Exit(1)
		}
	} else {
		flag.Usage()
	}
	os.Exit(0)
}
