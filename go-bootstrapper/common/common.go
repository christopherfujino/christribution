package common

import (
	"fmt"
	"os"
)

const Ignore = "./ignore"

var Archives = fmt.Sprintf("%s/archives", Ignore)
var Sources = fmt.Sprintf("%s/sources", Ignore)
const ManifestPath = "./manifest.json"

func EnsureDir(path string) {
	var _, err = os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0755)
		if err != nil {
			panic(err)
		}
	}
}

type Archive struct {
	Remote string
}
