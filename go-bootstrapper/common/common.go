package common

import (
	"fmt"
	"os"
	"path"
	"time"
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

type Manifest struct {
	Date     time.Time `json:"date"`
	Archives []Archive `json:"archives"`
}

func CreateArchive(name string, remote string, hash string) Archive {
	return Archive{
		Name:      name,
		Remote:    remote,
		Hash:      hash,
		LocalPath: fmt.Sprintf("%s/%s", Archives, path.Base(remote)),
	}
}

type Archive struct {
	Name      string `json:"name"`
	Remote    string `json:"remote"`
	Hash      string `json:"hash"`
	LocalPath string `json:"localPath"`
}
