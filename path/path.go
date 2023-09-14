package path

import (
	"errors"
	"fmt"
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
	file_identifier string // 原始字符串，注册时候字面量书写的
	file            string
	dir             string
	isReg           bool
	reg             *regexp.Regexp
}

func is_reg(filename string) (b bool, realFilename string) {
	realFilename = filename
	if b = strings.HasPrefix(filename, "reg:"); b {
		realFilename = realFilename[4:]
	}
	return
}

// if reg filename startwith reg:
func NewPath(file_identifier string) (*Path, error) {
	if file_identifier == "" {
		return nil, errors.New("file_identifier is empty")
	}
	// 最后是有斜杠的
	dir, file := split(file_identifier)
	if file == "" {
		return nil, errors.New("file is empty")
	}
	b, realFilename := is_reg(file)

	c := &Path{
		file_identifier: file_identifier,
		file:            realFilename,
		dir:             dir,
		isReg:           b,
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
	if this.file_identifier == delta {
		return true, ""
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

func (this *Path) IsReg() bool {
	return this.isReg
}

func (this *Path) Dir() string {
	return this.dir
}

func (this *Path) FileIdentifier() string {
	return this.file_identifier
}

func EqualDir(dir1, dir2 string) bool {
	if dir1 == dir2 {
		return true
	}
	return abs_dir(dir1) == abs_dir(dir2)
}

// join 本身没有斜杠
func abs_dir(dir string) string {
	if !filepath.IsAbs(dir) {
		return fmt.Sprintf("%s%c", filepath.Join(Pwd, dir), filepath.Separator)
	} else {
		if !strings.HasSuffix(dir, "/") {
			return fmt.Sprintf("%s%c", dir, filepath.Separator)
		}
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
