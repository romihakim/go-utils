package utils

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	trealpath1, _ := RealPath("/home/go/../go/test/../")
	equal(t, "D:\\home\\go", trealpath1)

	tbasename := BaseName("/home/go/src/pkg/utils.go")
	equal(t, "utils.go", tbasename)

	tPathinfo := PathInfo("/home/go/utils.go.go", -1)
	equal(t, map[string]string{"dirname": "\\home\\go", "basename": "utils.go.go", "extension": "go", "filename": "utils.go"}, tPathinfo)

	// tDiskFreeSpace, _ := DiskFreeSpace("/")
	// gt(t, float64(tDiskFreeSpace), 0)

	// tDiskTotalSpace, _ := DiskTotalSpace("/")
	// gte(t, float64(tDiskTotalSpace), 0)

	wd, _ := os.Getwd()
	tfilesize, _ := FileSize(wd)
	gt(t, float64(tfilesize), 0)
}
