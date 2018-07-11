package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-courier/loaderx"
)

func TestGenerator(t *testing.T) {
	cwd, _ := os.Getwd()
	p, pkgInfo, _ := loaderx.LoadWithTests(filepath.Join(cwd, "../__examples__"))

	g := NewStatusErrorGenerator(p, pkgInfo)

	g.Scan("StatusError")
	g.Output(cwd)
}
