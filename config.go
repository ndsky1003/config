package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/fsnotify/fsnotify"

	"github.com/ndsky1003/config/options"
)

type LoadFunc func([]byte) (any, error)

type config_mgr struct {
	items     []*load_item
	opt       *options.Option
	reWatch   chan struct{}
	startFlag uint32
}

func New(opts ...*options.Option) *config_mgr {
	opt := options.New().Merge(opts...)
	// 默认值
	if opt.Dirs == nil {
		opt.Dirs = []string{}
	}
	if opt.AutoLoad == nil {
		a := true
		opt.AutoLoad = &a
	}
	// 默认值
	c := &config_mgr{
		opt:     opt,
		reWatch: make(chan struct{}),
	}
	if *opt.AutoLoad {
		go c.AutoLoad()
	}
	return c
}

type load_item struct {
	rt     reflect.Type
	rv     unsafe.Pointer
	file   string
	isReg  bool
	reg    regexp.Regexp
	dirs   []string
	f      LoadFunc
	isAuto bool
}

type load_reg_item struct {
	rt     reflect.Type
	rv     unsafe.Pointer
	file   string
	isReg  bool
	reg    regexp.Regexp
	dirs   []string
	f      LoadFunc
	isAuto bool
}

func (this *load_item) equal(opt *load_item) bool {
	if this == nil || opt == nil {
		return false
	}
	if this.rt == opt.rt {
		return true
	}
	return false
}

func isRegexp(filename string) bool {
	return strings.HasPrefix(filename, "/") && strings.HasSuffix(filename, "/")
}

var default_config_mgr = New()

func Regist[T any](filename string, fn LoadFunc, opts ...*options.Option) error {
	dir, file := filepath.Split(filename)
	if file == "" {
		return errors.New("file is invalid")
	}
	opt := options.New().Merge(opts...)
	if dir != "" {
		opt.SetDir(dir)
	}
	if len(opt.Dirs) == 0 {
		opt.SetDir(default_config_mgr.opt.Dirs...)
	}

	var a T
	var rv unsafe.Pointer
	rt := reflect.TypeOf(a)
	switch rt.Kind() {
	case reflect.Map:
		m := reflect.MakeMap(rt).Interface().(T)
		rv = unsafe.Pointer(&m)
	case reflect.Slice:
		m := reflect.MakeSlice(rt, 0, 0).Interface().(T)
		rv = unsafe.Pointer(&m)
	case reflect.Chan:
		m := reflect.MakeChan(rt, 0).Interface().(T)
		rv = unsafe.Pointer(&m)
	case reflect.Pointer:
		return errors.New("regist must not pointer")
	case reflect.Func:
		return errors.New("regist must not func")
	case reflect.Interface:
		return errors.New("regist must not Interface")
	default:
		rv = unsafe.Pointer(&a)
	}
	item := &load_item{
		rt:   rt,
		rv:   rv,
		file: file,
		dirs: opt.Dirs,
		f:    fn,
	}
	default_config_mgr.items = append(default_config_mgr.items, item)
	return nil
}

// 保证*T不为nil
func Get[T any]() *T {
	var a T
	rt := reflect.TypeOf(a)
	for _, item := range default_config_mgr.items {
		if item.rt == rt {
			return (*T)(atomic.LoadPointer(&item.rv))
		}
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

func (this *config_mgr) loadfile(pathname string) error {
	// d := m.find(file)
	// if d == nil {
	// 	return
	// }

	// buf, err := os.ReadFile(filename)
	// if err != nil {
	// 	return err
	// }

	// md5sum := md5.Sum(buf)
	// newMD5 := fmt.Sprintf("%x", md5sum)

	// if newMD5 == d.md5 {
	// 	logger.Infof("%v 没有变动, 不需要加载\n", file)
	// 	return
	// } else {
	// 	logger.Infof("开始加载：%v\n", file)
	// }
	//
	// err = d.f(buf)
	//
	// if err != nil {
	// 	logger.Err("加载失败", file, err)
	// } else {
	// 	logger.Infof("配置%v 加载成功, md5: %v\n", file, newMD5)
	// 	d.md5 = newMD5
	// }
	//
	return nil
}

func (this *config_mgr) AutoLoad() {
	if b := atomic.CompareAndSwapUint32(&(this.startFlag), 0, 1); !b {
		return
	}
	defer func() {
		fmt.Println("AutoLoad is done")
	}()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()
here:
	for _, workdir := range this.opt.Dirs {
		if _, err = os.Stat(workdir); os.IsNotExist(err) {
			if err = os.MkdirAll(workdir, 0777); err != nil {
				fmt.Println("err:", err)
			}
		}

		if err = watcher.Add(workdir); err != nil {
			fmt.Println("add dir err:", err)
		}
		fmt.Println("ddd:", workdir)
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("modified file:", event.Name, filepath.IsAbs(event.Name))
			}
		case err := <-watcher.Errors:
			fmt.Println("error:", err)
		case <-this.reWatch:
			goto here
		}
	}
}
