package item

import (
	"crypto/md5"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"unsafe"

	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
)

type (
	LoadFunc[T any]    func([]byte) (*T, error)
	LoadRegFunc[T any] func([]string, []byte) (*T, error)
)

type IItem[T any] interface {
	Match(pathname string) (bool, string)
	LoadFile(pathname string, buf []byte) error
	Path() *path.Path
	RT() reflect.Type
	RVS() []*ItemValue
	CheckBuf([]byte) error // buf检查
	Opts() *options.Option[T]
}

type Item[T any] struct {
	T     reflect.Type
	VS    []*ItemValue
	F     LoadFunc[T]    //这个也支持正则,只是不支持加载函数无法探测处flag
	F_reg LoadRegFunc[T] //与 F,二选一
	P     *path.Path
	Opt   *options.Option[T]
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
	// _, err := this.F(buf)
	_, err := this.call_F([]string{}, buf)
	return err
}

func (this *Item[T]) Opts() *options.Option[T] {
	return this.Opt
}

func (this *Item[T]) call_F(submath []string, buf []byte) (*T, error) {
	if this.F != nil {
		return this.F(buf)
	} else if this.F_reg != nil {
		return this.F_reg(submath, buf)
	}
	return nil, errors.New("no load func")

}

func (this *Item[T]) LoadFile(file_identifier string, buf []byte) error {
	filename := filepath.Base(file_identifier)
	submatch, flag := this.P.Flag(filename)
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
	// pv, err := this.F(buf)
	pv, err := this.call_F(submatch, buf)
	if err != nil {
		fmt.Printf("加载失败:%v,%v\n", file_identifier, err)
		return err
	}
	if opt := this.Opts(); opt != nil && opt.SuccessFunc != nil {
		(*opt.SuccessFunc)(pv)
	}
	fmt.Printf("配置%v 加载成功, md5: %v\n", file_identifier, newMD5)
	if need_append {
		this.VS = append(this.VS, tmp_rv_meta)
	}
	tmp_rv_meta.V = reflect.ValueOf(pv).UnsafePointer()
	tmp_rv_meta.MD5 = newMD5
	return nil
}
