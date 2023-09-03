package path

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path, err := NewPath("/Users/ppll/go/workSpace/self-pkg/config/path/path.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(path.Match("./path.go"))
}

func TestABSDir(t *testing.T) {
	// fmt.Println(abs_dir("/Users/ppll/go/workSpace/self-pkg/config/path/"))
	fmt.Println(filepath.Join(Pwd, "../path", "/"))
	fmt.Println(abs_dir("."))
	fmt.Println(abs_dir("./"))
}

func TestSplite(t *testing.T) {
	fmt.Println(split("../cc/path.go"))
}
