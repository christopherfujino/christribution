package fetch

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

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

	batchDownload(remotesToFetch(manifest.Archives))

	batchDownload(remotesToFetch(manifest.Patches))
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

func batchDownload(archives []common.Archive) {
	var wg sync.WaitGroup
	var archiveChan = make(chan common.Archive, 100)
	var worker = func(id int, archiveChan <-chan common.Archive, wg *sync.WaitGroup) {
		defer wg.Done()
		for archive := range archiveChan {
			fmt.Printf("[%d] Starting job for %s...\n", id, archive.Remote)
			downloadAndVerify(archive)
			fmt.Printf("[%d] Finished job for %s...\n", id, archive.Remote)
		}
	}

	for i := 0; i < common.CONCURRENT_WORKERS; i += 1 {
		wg.Add(1)
		go worker(i, archiveChan, &wg)
	}

	for _, archive := range archives {
		fmt.Printf("Sending job for %s to worker pool...\n", archive.Remote)
		archiveChan <- archive
	}

	close(archiveChan)
	wg.Wait()
}

func downloadAndVerify(archive common.Archive) {
	var success = download(archive)
	if success {
		verify(archive)
		fmt.Printf("Download of %s successful.\n", archive.LocalPath)
	}
}

func download(archive common.Archive) (isDone bool) {
	defer (func() {
		if !isDone {
			var err = os.Remove(archive.LocalPath)
			if err == nil {
				fmt.Printf("Removed %s after failed CURL\n", archive.LocalPath)
			}
		}
	})()

	fmt.Printf("GET %s\n", archive.Remote)
	cmd := exec.Command("curl", "-L", archive.Remote, "-o", archive.LocalPath)
	var buffer = strings.Builder{}
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	var err = cmd.Run()
	if err != nil {
		_ = os.Remove(archive.LocalPath)
		fmt.Fprintf(os.Stderr, "  -> [ERROR] %v\n", err)
		fmt.Fprintf(os.Stderr, "\n%s\n", buffer.String())
		// No signal that we failed, user will have to re-run the fetch to ensure
		// there is no remaining work to be done.
		return
	}

	return true
}

func verify(archive common.Archive) {
	// check hash
	var hash = md5.New()
	localFile, err := os.Open(archive.LocalPath)
	if err != nil {
		panic(err)
	}
	defer localFile.Close()
	_, err = io.Copy(hash, localFile)
	if err != nil {
		panic(err)
	}
	var hashBytes = hash.Sum(nil)
	var hashString = hex.EncodeToString(hashBytes)
	if hashString != archive.Hash {
		panic(fmt.Sprintf(
			"Expected %s to have a hash of %s but it actually had %s",
			archive.Name,
			archive.Hash,
			hashString,
		))
	}
}
