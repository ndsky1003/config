package path

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Pwd string

func init() {
	var err error
	Pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

type Path struct {
	file  string
	dir   string
	isReg bool
	reg   *regexp.Regexp
}

func is_reg(filename string) (b bool, realFilename string) {
	realFilename = filename
	if b = strings.HasPrefix(filename, "reg:"); b {
		realFilename = realFilename[4:]
	}
	return
}

// if reg filename startwith reg:
func NewPath(filename string) (*Path, error) {
	if filename == "" {
		return nil, errors.New("filename is empty")
	}
	// 最后是有斜杠的
	dir, file := split(filename)
	if file == "" {
		return nil, errors.New("file is empty")
	}
	b, realFilename := is_reg(file)

	c := &Path{
		file:  realFilename,
		dir:   dir,
		isReg: b,
	}
	if b {
		c.reg = regexp.MustCompile(realFilename)
	}
	return c, nil
}

// delta 是文件的路径
// return
// bool 表示是否匹配
// string 表示匹配上的文件名，用于提取标示符
func (this *Path) Match(delta string) (bool, string) {
	if delta == "" {
		return false, ""
	}
	tmpdir, tmpfile := split(delta)
	selfdir := abs_dir(this.dir)
	tmpdir = abs_dir(tmpdir)
	if selfdir != tmpdir {
		return false, ""
	}
	if this.isReg {
		return this.reg.MatchString(tmpfile), tmpfile
	} else {
		return this.file == tmpfile, tmpfile
	}
}

func (this *Path) Flag(filename string) string {
	if this.isReg {
		return strings.Join(this.reg.FindStringSubmatch(filename)[1:], "_")
	} else {
		return ""
	}
}

func (this *Path) File() string {
	return this.file
}

func (this *Path) Dir() string {
	return this.dir
}

func (this *Path) List() ([]string, error) {
	if !this.isReg {
		return []string{filepath.Join(this.Dir(), this.File())}, nil
	}
	dir := this.Dir()
here:
	file, err := os.Open(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		} else {
			goto here
		}
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0)
	for _, file := range names {
		realPath := filepath.Join(dir, file)
		if b, _ := this.Match(realPath); b {
			res = append(res, realPath)
		}
	}
	return res, nil
}

// join 本身没有斜杠
func abs_dir(dir string) string {
	if !filepath.IsAbs(dir) {
		return filepath.Join(Pwd, dir) + "/"
	}
	return dir
}

// 保证dir后面都有分隔符
func split(path string) (dir, file string) {
	if path == "" {
		return
	}
	dir, file = filepath.Split(filepath.Clean(path))
	if dir == "" {
		dir = "./"
	}
	return
}
