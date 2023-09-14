package path

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path, err := NewPath("reg:db_vip_([a-z]{3})_([a-z]{3}).yaml")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(path.Match("db_vip_bob_dod.yaml"))
	fmt.Println(path.Flag("db_vip_bob_dod.yaml"), "dddd")
}

func TestABSDir(t *testing.T) {
	// fmt.Println(abs_dir("/Users/ppll/go/workSpace/self-pkg/config/path/"))
	fmt.Println(filepath.Join(Pwd, "../path", "/"))
	fmt.Println(abs_dir("."))
	fmt.Println(abs_dir("./"))
}

func TestSplite(t *testing.T) {
	pwd, err := os.Getwd()
	t.Error(err)
	t.Log(pwd)
	t.Log(filepath.IsAbs(pwd))
	t.Log(EqualDir("./", "."))
	t.Log(EqualDir("./", pwd))
	t.Log(EqualDir(pwd, "."))
}
