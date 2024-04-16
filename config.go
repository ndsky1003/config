// 该文件的意义，纯粹的就是因为go不支持方法泛型的一个包装
package config

import (
	"reflect"
	"sync/atomic"

	"github.com/ndsky1003/config/checker"
	"github.com/ndsky1003/config/item"
	"github.com/ndsky1003/config/options"
	"github.com/ndsky1003/config/path"
	"github.com/ndsky1003/config/watcher"
)

func Stop() {
	default_config_mgr.Stop()
}

func SetChecker(c checker.IChecker) {
	default_config_mgr.SetChecker(c)
}

func SetCheckerIdentifierFunc(f func(string) string) {
	default_config_mgr.SetCheckerIdentifierFunc(f)
}

func SetWatcher(w watcher.IWatcher) {
	default_config_mgr.SetWatcher(w)
}

func Regist[T any](file_identifier string, fn item.LoadFunc[T], opts ...*options.Option) error {
	path, err := path.NewPath(file_identifier)
	if err != nil {
		return err
	}
	opt := options.New().Merge(opts...)
	var a T
	rt := reflect.TypeOf(a)
	item := &item.Item[T]{
		T:   rt,
		F:   fn,
		P:   path,
		Opt: opt,
	}
	return default_config_mgr.Regist(item)
}

// 保证*T不为nil ,找不到就是零值，且map、slice、chan的零值是初始化过的
// flag 最多只有一个值，不存在的时候，获取的是非正则匹配的时候
// flag 只有一个值的时候，是正则匹配的提取物的下划线作为连接符号的标志
// 可变参数，纯粹为了美观
// 返回值必须指针，否则失去其意义
func Get[T any](flags ...string) *T {
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
		return (*T)(atomic.LoadPointer(&v.V))
	}
	switch rt.Kind() {
	case reflect.Map:
		m := reflect.MakeMap(rt).Interface().(T)
		return &m
	case reflect.Slice:
		m := reflect.MakeSlice(rt, 0, 0).Interface().(T)
		return &m
	case reflect.Chan:
		m := reflect.MakeChan(rt, 0).Interface().(T)
		return &m
	}
	return &a
}

// 针对引用类型，不返回其指针
// 专门针对map,slice,chan使用
func GetRef[T any](flags ...string) T {
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
		return *(*T)(atomic.LoadPointer(&v.V))
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
