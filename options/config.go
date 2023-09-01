package options

import (
	"os"
	"path/filepath"

	"github.com/samber/lo"
)

var pwd string

func init() {
	var err error
	if pwd, err = os.Getwd(); err != nil {
		panic(err)
	}
}

type Option struct {
	Dirs     []string
	AutoLoad *bool
}

func New() *Option {
	return &Option{}
}

func (this *Option) SetDir(dirs ...string) *Option {
	for _, dir := range dirs {
		if dir != "" {
			if !filepath.IsAbs(dir) {
				dir = filepath.Join(pwd, dir)
			}
			if !lo.Contains(this.Dirs, dir) {
				this.Dirs = append(this.Dirs, dir)
			}
		}
	}
	return this
}

func (this *Option) SetAutoLoad(b bool) *Option {
	this.AutoLoad = &b
	return this
}

func (this *Option) Merge(opts ...*Option) *Option {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

func (this *Option) merge(opt *Option) {
	for _, dir := range opt.Dirs {
		if dir != "" && !lo.Contains(this.Dirs, dir) {
			this.Dirs = append(this.Dirs, dir)
		}
	}
	if opt.AutoLoad != nil {
		this.AutoLoad = opt.AutoLoad
	}
}
