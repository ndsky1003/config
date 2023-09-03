package config

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
)

type i_load_item interface {
	Match(pathname string) (bool, string)
	LoadFile(pathname string) error
	RT() reflect.Type
	RVS() []*load_item_meta
}

type load_item[T any] struct {
	rt   reflect.Type
	rvs  []*load_item_meta
	f    LoadFunc[T]
	path *path.Path
	opt  *options.Option
}

type load_item_meta struct {
	flag string
	rv   unsafe.Pointer
	md5  string
}

func (this *load_item[T]) RT() reflect.Type {
	return this.rt
}

func (this *load_item[T]) RVS() []*load_item_meta {
	return this.rvs
}

func (this *load_item[T]) Match(pathname string) (bool, string) {
	return this.path.Match(pathname)
}

func (this *load_item[T]) load_files() error {
	files, err := this.path.List()
	if err != nil {
		return err
	}
	for _, pathname := range files {
		if err = this.LoadFile(pathname); err != nil {
			return err
		}
	}
	return nil
}

func (this *load_item[T]) LoadFile(pathname string) error {
	buf, err := os.ReadFile(pathname)
	if err != nil {
		return err
	}
	filename := filepath.Base(pathname)
	flag := this.path.Flag(filename)
	md5sum := md5.Sum(buf)
	newMD5 := fmt.Sprintf("%x", md5sum)

	var tmp_rv_meta *load_item_meta
	for _, rv_meta := range this.rvs {
		if rv_meta.flag == flag {
			tmp_rv_meta = rv_meta
		}
	}
	var need_append bool
	if tmp_rv_meta == nil {
		tmp_rv_meta = &load_item_meta{
			flag: flag,
		}
		need_append = true
	}
	if tmp_rv_meta.md5 == newMD5 {
		fmt.Printf("%v 没有变动, 不需要加载\n", pathname)
		return nil
	}
	fmt.Printf("开始加载：%v\n", pathname)
	pv, err := this.f(buf)
	if err != nil {
		fmt.Printf("加载失败:%v,%v\n", pathname, err)
		return err
	}
	fmt.Printf("配置%v 加载成功, md5: %v\n", pathname, newMD5)
	if need_append {
		this.rvs = append(this.rvs, tmp_rv_meta)
	}
	tmp_rv_meta.rv = reflect.ValueOf(pv).UnsafePointer()
	tmp_rv_meta.md5 = newMD5
	return nil
}
