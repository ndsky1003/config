package watcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// 默认文件监视器
type file_watcher struct {
	Dir         string
	done        chan struct{}
	dispense_fn func(file_identifier string, buf []byte) error
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

func (this *file_watcher) SetReloadData(f func(file_identifier string, buf []byte) error) error {
	this.dispense_fn = f
	return nil
}

func (this *file_watcher) LoadFile(file_identifier string) error {
	buf, err := os.ReadFile(file_identifier)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return errors.New("buf is nil")
	}
	return this.dispense_fn(file_identifier, buf)
}

func (this *file_watcher) LoadFiles(fn func(file_identifier string) bool) error {
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
			if err := this.LoadFile(realPath); err != nil {
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
				if err := this.LoadFile(event.Name); err != nil {
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
