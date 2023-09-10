// 该文件的意义，纯粹的就是因为go不支持方法泛型的一个包装
package config

import (
	"reflect"
	"sync/atomic"

	"github.com/ndsky1003/config/checker"
	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
	"github.com/ndsky1003/config/watcher"
)

type (
	LoadFunc[T any] func([]byte) (*T, error)
)

func Stop() {
	default_config_mgr.Stop()
}

func SetChecker(c checker.IChecker) {
	default_config_mgr.SetChecker(c)
}

func SetWatcher(w watcher.IWatcher) {
	default_config_mgr.SetWatcher(w)
}

func Regist[T any](file_identifier string, fn LoadFunc[T], opts ...*options.Option) error {
	path, err := path.NewPath(file_identifier)
	if err != nil {
		return err
	}
	opt := options.New().Merge(opts...)
	var a T
	rt := reflect.TypeOf(a)
	item := &load_item[T]{
		rt:   rt,
		f:    fn,
		path: path,
		opt:  opt,
	}
	return default_config_mgr.RegistLoadItem(item)
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
	if v := default_config_mgr.GetLoadItem(rt, flag); v != nil {
		return *(*T)(atomic.LoadPointer(&v.rv))
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
