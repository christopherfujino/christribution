package extract

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/christopherfujino/christribution/go-bootstrapper/common"
)

func Extract() {
	entries, err := os.ReadDir(common.Archives)
	if err != nil {
		panic(err)
	}
	common.EnsureDir(common.Sources)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		var baseName = entry.Name()
		var fullName = fmt.Sprintf("%s/%s", common.Archives, baseName)
		// Not all tar archives have an inner directory
		var sourceNamespace = fmt.Sprintf("%s/%s", common.Sources, baseName)
		common.EnsureDir(sourceNamespace)

		if strings.HasSuffix(baseName, "tar.xz") {
			run("tar", "xvf", fullName, "-C", sourceNamespace)
		} else if strings.HasSuffix(baseName, ".tar.gz") {
			run("tar", "xvf", fullName, "-C", sourceNamespace)
		} else if strings.HasSuffix(baseName, ".tar.bz2") {
			run("tar", "xvf", fullName, "-C", sourceNamespace)
		} else if strings.HasSuffix(baseName, ".tgz") {
			run("tar", "xvf", fullName, "-C", sourceNamespace)
		} else {
			panic(baseName)
		}
	}
}

func run(args ...string) {
	var first = args[0]
	var rest = args[1:]
	var cmd = exec.Command(first, rest...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	var err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
