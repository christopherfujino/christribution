package fetch

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/christopherfujino/christribution/go-bootstrapper/common"
)

func Fetch() {
	_, err := os.Stat(common.Ignore)
	if err != nil {
		// Needs execute permissions to access
		os.Mkdir(common.Ignore, 0755)
	}

	_, err = os.Stat(common.Archives)
	if err != nil {
		os.Mkdir(common.Archives, 0755)
	}

	manifestBytes, err := os.ReadFile(common.ManifestPath)
	if err != nil {
		panic(err)
	}
	var manifest common.Manifest
	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		panic(err)
	}
	var remotes = remotesToFetch(manifest.Archives)

	download(remotes)
}

func remotesToFetch(archives []common.Archive) []common.Archive {
	var outputArchives []common.Archive

	for _, archive := range archives {
		if archive.LocalPath == "" {
			panic(fmt.Sprintf("Found an empty LocalPath in %v", archive))
		}
		_, err := os.Stat(archive.LocalPath)
		if err != nil {
			outputArchives = append(outputArchives, archive)
		} else {
			fmt.Printf("The file %s already exists, skipping fetch.\n", archive.LocalPath)
		}
	}

	return outputArchives
}

func download(archives []common.Archive) {
	for _, archive := range archives {
		var remotePath = archive.Remote
		var localPath = archive.LocalPath
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
			panic(
				fmt.Sprintf(
					"Error creating the local file %s:\n%s",
					localPath,
					err.Error(),
				),
			)
		}
		defer writeFile.Close()

		{
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
		}

		// check hash
		{
			var hash = md5.New()
			localFile, err := os.Open(archive.LocalPath)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(hash, localFile)
			if err != nil {
				panic(err)
			}
			var hashBytes = hash.Sum(nil)
			var hashString = hex.EncodeToString(hashBytes)
			if hashString != archive.Hash {
				panic(
					fmt.Sprintf(
						"Expected %s to have a hash of %s but it actually had %s",
						archive.Name,
						archive.Hash,
						hashString,
					),
				)
			}
		}

		isDone = true
		fmt.Printf("Download of %s successful.\n", localPath)
	}
}
