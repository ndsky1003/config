package watcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/ndsky1003/config/item"
)

// 默认文件监视器
type file_watcher struct {
	Dir             string
	done            chan struct{}
	distribute_func func(file_identifier string, buf []byte) error
}

func NewFileWatcher(dir string) (*file_watcher, error) {
	if dir == "" {
		return nil, errors.New("dis not empty")
	}
	file, err := os.Open(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
	}
	file.Close()
	c := &file_watcher{
		Dir:  dir,
		done: make(chan struct{}),
	}
	go c.protect_run()
	return c, nil
}

func (this *file_watcher) Stop() error {
	if this.done != nil {
		close(this.done)
	}
	this.done = nil
	return nil
}

func (this *file_watcher) SetDistributeFunc(
	f func(file_identifier string, buf []byte) error,
) error {
	this.distribute_func = f
	return nil
}

func (this *file_watcher) Regist(item item.IItem) (err error) {
	if item.Path().IsReg() {
		err = this.load_files(func(file_identifier1 string) bool {
			b, _ := item.Match(file_identifier1)
			return b
		})
	} else {
		err = this.load_file(item.Path().FileIdentifier())
	}
	return
}

func (this *file_watcher) load_file(file_identifier string) error {
	buf, err := os.ReadFile(file_identifier)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return errors.New("buf is nil")
	}
	return this.distribute_func(file_identifier, buf)
}

func (this *file_watcher) load_files(fn func(file_identifier string) bool) error {
	dir := this.Dir
	file, err := os.Open(dir)
	if err != nil {
		return err
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, file := range names {
		realPath := filepath.Join(dir, file)
		if fn(realPath) {
			if err := this.load_file(realPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *file_watcher) protect_run() {
	defer fmt.Println("config_mgr exit")
	for {
		this.auto_load()
		select {
		case <-this.done:
			break
		default:
		}
	}
}

func (this *file_watcher) auto_load() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	if err := watcher.Add(this.Dir); err != nil {
		panic(err)
	}
	defer watcher.Close()
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("loadFile:", event.Name)
				if err := this.load_file(event.Name); err != nil {
					fmt.Printf("loadfile:err:%v\n", err)
				}
			}
		case err := <-watcher.Errors:
			panic(err)
		case <-this.done:
			return
		}
	}
}
