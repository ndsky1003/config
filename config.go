package config

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
)

type (
	LoadFunc[T any] func([]byte) (*T, error)
)

func Stop() {
	default_config_mgr.stop()
}

func Regist[T any](filename string, fn LoadFunc[T], opts ...*options.Option) error {
	path, err := path.NewPath(filename)
	if err != nil {
		return err
	}
	default_config_mgr.push_dir(path.Dir())
	opt := options.New().Merge(opts...)
	var a T
	rt := reflect.TypeOf(a)
	for _, item := range default_config_mgr.items {
		if item.RT() == rt {
			return fmt.Errorf("%v exist,please rename", rt.String())
		}
	}
	item := &load_item[T]{
		rt:   rt,
		f:    fn,
		path: path,
		opt:  opt,
	}
	if err = item.load_files(); err != nil {
		return err
	}
	default_config_mgr.items = append(default_config_mgr.items, item)
	return nil
}

// 保证*T不为nil ,找不到就是零值，且map、slice、chan的零值是初始化过的
// flag 最多只有一个值，不存在的时候，获取的是非正则匹配的时候
// flag 只有一个值的时候，是正则匹配的提取物的下划线作为连接符号的标志
// 可变参数，纯粹为了美观
func Get[T any](flags ...string) T {
	if len(flags) >= 2 {
		panic("params  most one")
	}
	var flag string
	if len(flags) == 1 {
		flag = flags[0]
	}
	var a T
	rt := reflect.TypeOf(a)
	for _, item := range default_config_mgr.items {
		if item.RT() == rt {
			for _, item_meta := range item.RVS() {
				if flag == item_meta.flag {
					return *(*T)(atomic.LoadPointer(&item_meta.rv))
				}
			}
		}
	}
	switch rt.Kind() {
	case reflect.Map:
		m := reflect.MakeMap(rt).Interface().(T)
		return m
	case reflect.Slice:
		m := reflect.MakeSlice(rt, 0, 0).Interface().(T)
		return m
	case reflect.Chan:
		m := reflect.MakeChan(rt, 0).Interface().(T)
		return m
	}
	return a
}
