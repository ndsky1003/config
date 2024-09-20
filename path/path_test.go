package path

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path, err := New("reg:db_vip_([a-z]{3})_([a-z]{3}).yaml", true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(path.Match("db_vip_bob_dod.yaml"))
	v, f := path.Flag("db_vip_bob_dod.yaml")
	fmt.Println(v, f)
}

func TestABSDir(t *testing.T) {
	// fmt.Println(abs_dir("/Users/ppll/go/workSpace/self-pkg/config/path/"))
	fmt.Println(filepath.Join(pwd, "../path", "/"))
	fmt.Println(abs_dir("."))
	fmt.Println(abs_dir("./"))
}

func TestSplite(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	t.Log(pwd)
	t.Log(filepath.IsAbs(pwd))
	t.Log(EqualDir("./", "."))
	t.Log(EqualDir("./", pwd))
	t.Log(EqualDir(pwd, "."))
}
