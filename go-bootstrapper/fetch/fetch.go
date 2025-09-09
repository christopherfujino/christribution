package fetch

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/christopherfujino/christribution/go-bootstrapper/common"
)

func Fetch() {
	_, err := os.Stat(common.Ignore)
	if err != nil {
		// Needs execute to access
		os.Mkdir(common.Ignore, 0755)
	}

	_, err = os.Stat(common.Archives)
	if err != nil {
		os.Mkdir(common.Archives, 0755)
	}

	manifestBytes, err := os.ReadFile("./manifest.json")
	if err != nil {
		panic(err)
	}
	var remotePaths []string
	err = json.Unmarshal(manifestBytes, &remotePaths)
	if err != nil {
		panic(err)
	}
	var remotes = remotesToFetch(remotePaths)
	//fmt.Println(remotes)
	download(remotes)
}

func remotesToFetch(remotePaths []string) [][]string {
	var remotes [][]string

	for _, remotePath := range remotePaths {
		var localPath = fmt.Sprintf("%s/%s", common.Ignore, path.Base(remotePath))
		_, err := os.Stat(localPath)
		if err != nil {
			remotes = append(remotes, []string{remotePath, localPath})
		} else {
			fmt.Printf("The file %s already exists, skipping fetch.\n", localPath)
		}
	}

	return remotes
}

func download(remotes [][]string) {
	for _, tuple := range remotes {
		var remotePath = tuple[0]
		var localPath = tuple[1]
		var isDone = false

		defer (func() {
			if !isDone {
				fmt.Printf("Removing %s...\n", localPath)
				err := os.Remove(localPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
				}
			}
		})()

		writeFile, err := os.Create(localPath)
		if err != nil {
			panic(err)
		}
		defer writeFile.Close()

		fmt.Printf("GET %s\n", remotePath)
		cmd := exec.Command("curl", "-L", remotePath, "-o", localPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			_ = os.Remove(localPath)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		isDone = true
		fmt.Printf("Download of %s successful.\n", localPath)
	}
}
