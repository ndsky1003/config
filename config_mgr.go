package config

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

var default_config_mgr = new()

type config_mgr struct {
	done        chan struct{}
	items       []i_load_item
	dirs        []string
	watcher_dir chan string
}

func new() *config_mgr {
	c := &config_mgr{
		watcher_dir: make(chan string),
	}
	go c.run()
	return c
}

func (this *config_mgr) push_dir(dir string) {
	if !lo.Contains(this.dirs, dir) {
		this.watcher_dir <- dir
		this.dirs = append(this.dirs, dir)
	}
}

func (this *config_mgr) loadfile(pathname string) {
	fmt.Println("loadfile:", pathname)
	for _, item := range this.items {
		fmt.Println(item)
		if b, _ := item.Match(pathname); b {
			fmt.Println("math:", pathname)
			if err := item.LoadFile(pathname); err != nil {
				continue
			}
		}
	}
}

// 并非线程安全，没必要
func (this *config_mgr) stop() {
	close(this.done)
	this.done = nil
}

func (this *config_mgr) run() {
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

func (this *config_mgr) auto_load() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err:", err)
		}
	}()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()
	for _, workdir := range this.dirs {
		if _, err := os.Stat(workdir); os.IsNotExist(err) {
			if err = os.MkdirAll(workdir, 0777); err != nil {
				fmt.Println("err:", err)
			}
		}
		if err := watcher.Add(workdir); err != nil {
			fmt.Println("add dir err:", err)
			continue
		}
		fmt.Println("add dir:", workdir)
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				this.loadfile(event.Name)
			}
		case err := <-watcher.Errors:
			fmt.Println("error:", err)
		case dir := <-this.watcher_dir:
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if err = os.MkdirAll(dir, 0777); err != nil {
					fmt.Println("err:", err)
				}
			}
			if err := watcher.Add(dir); err != nil {
				fmt.Println("add dir err:", err)
			}
			fmt.Println("add dir:", dir)
		case <-this.done:
			return
		}
	}
}
