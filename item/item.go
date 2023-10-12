package item

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
)

type (
	LoadFunc[T any] func([]byte) (*T, error)
)

type IItem interface {
	Match(pathname string) (bool, string)
	LoadFile(pathname string, buf []byte) error
	Path() *path.Path
	RT() reflect.Type
	RVS() []*ItemValue
	CheckBuf([]byte) error // buf检查
}

type Item[T any] struct {
	T   reflect.Type
	VS  []*ItemValue
	F   LoadFunc[T]
	P   *path.Path
	Opt *options.Option
}

type ItemValue struct {
	Flag string
	V    unsafe.Pointer
	MD5  string
}

func (this *Item[T]) RT() reflect.Type {
	return this.T
}

func (this *Item[T]) RVS() []*ItemValue {
	return this.VS
}

func (this *Item[T]) Match(pathname string) (bool, string) {
	return this.P.Match(pathname)
}

func (this *Item[T]) Path() *path.Path {
	return this.P
}

func (this *Item[T]) CheckBuf(buf []byte) error {
	_, err := this.F(buf)
	return err
}

func (this *Item[T]) LoadFile(file_identifier string, buf []byte) error {
	filename := filepath.Base(file_identifier)
	flag := this.P.Flag(filename)
	md5sum := md5.Sum(buf)
	newMD5 := fmt.Sprintf("%x", md5sum)

	var tmp_rv_meta *ItemValue
	for _, rv_meta := range this.VS {
		if rv_meta.Flag == flag {
			tmp_rv_meta = rv_meta
		}
	}
	var need_append bool
	if tmp_rv_meta == nil {
		tmp_rv_meta = &ItemValue{
			Flag: flag,
		}
		need_append = true
	}
	if tmp_rv_meta.MD5 == newMD5 {
		fmt.Printf("%v 没有变动, 不需要加载\n", file_identifier)
		return nil
	}
	fmt.Printf("开始加载：%v\n", file_identifier)
	pv, err := this.F(buf)
	if err != nil {
		fmt.Printf("加载失败:%v,%v\n", file_identifier, err)
		return err
	}
	fmt.Printf("配置%v 加载成功, md5: %v\n", file_identifier, newMD5)
	if need_append {
		this.VS = append(this.VS, tmp_rv_meta)
	}
	tmp_rv_meta.V = reflect.ValueOf(pv).UnsafePointer()
	tmp_rv_meta.MD5 = newMD5
	return nil
}
