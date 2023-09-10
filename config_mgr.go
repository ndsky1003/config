package config

import (
	"fmt"
	"reflect"

	"github.com/ndsky1003/config/checker"
	"github.com/ndsky1003/config/watcher"
)

const dir = "config"

var default_config_mgr *config_mgr

func init() {
	default_file_watcher, err := watcher.NewFileWatcher(dir)
	if err != nil {
		panic(err)
	}
	default_config_mgr = New(default_file_watcher)
}

type config_mgr struct {
	items   []i_load_item
	watcher watcher.IWatcher
	checker checker.IChecker
}

func New(watcher watcher.IWatcher) *config_mgr {
	c := &config_mgr{}
	c.SetWatcher(watcher)
	return c
}

func (this *config_mgr) Stop() {
	if this.watcher != nil {
		if err := this.watcher.Stop(); err != nil {
			fmt.Println("err:", err)
		}
	}
}

func (this *config_mgr) SetChecker(c checker.IChecker) {
	this.checker = c
}

func (this *config_mgr) SetWatcher(w watcher.IWatcher) {
	if this.watcher != nil {
		if err := this.watcher.Stop(); err != nil {
			fmt.Println("err:", err)
		}
	}

	this.watcher = w
	_ = this.watcher.SetReloadData(this.ReloadData)
}

func (this *config_mgr) RegistLoadItem(item i_load_item) error {
	file_identifier := item.Path().FileIdentifier()
	for _, v := range this.items {
		if v.RT() == item.RT() {
			return fmt.Errorf("%v exist,please rename", item.RT().String())
		}
		if file_identifier == v.Path().FileIdentifier() {
			return fmt.Errorf(
				"file_identifier:%v exist,please change",
				file_identifier,
			)
		}
	}

	this.items = append(this.items, item)
	var err error
	if item.Path().IsReg() {
		err = this.watcher.LoadFiles(func(file_identifier1 string) bool {
			b, _ := item.Match(file_identifier1)
			return b
		})
	} else {
		err = this.watcher.LoadFile(item.Path().FileIdentifier())
	}
	if err == nil && this.checker != nil {
		this.checker.On(file_identifier, this.ReloadData)
	}
	return err
}

func (this *config_mgr) GetLoadItem(rt reflect.Type, flag string) *load_item_meta {
	for _, item := range this.items {
		if item.RT() == rt {
			for _, item_meta := range item.RVS() {
				if flag == item_meta.flag {
					return item_meta
				}
			}
		}
	}
	return nil
}

// MARK check 也可注入该方法检查，成功才会替换
func (this *config_mgr) ReloadData(file_identifier string, buf []byte) error {
	for _, item := range this.items {
		if b, _ := item.Match(file_identifier); b {
			if err := item.LoadFile(file_identifier, buf); err != nil {
				continue
			}
		}
	}
	return nil
}
