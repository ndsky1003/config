package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ndsky1003/config/v2/checker"
	"github.com/ndsky1003/config/v2/item"
	"github.com/ndsky1003/config/v2/watcher"
)

var default_config_mgr *config_mgr

func init() {
	default_file_watcher, err := watcher.NewFileWatcher()
	if err != nil {
		panic(err)
	}
	default_config_mgr = New(default_file_watcher)
}

type config_mgr struct {
	items   []item.IItem
	watcher watcher.IWatcher
	checker checker.IChecker
	// gen_checker_identifier func(item.IItem) string
}

func New(watcher watcher.IWatcher) *config_mgr {
	c := &config_mgr{}
	if err := c.SetWatcher(watcher); err != nil {
		panic(err)
	}
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

// func (this *config_mgr) SetCheckerIdentifierFunc(f func(item.IItem) string) {
// 	this.gen_checker_identifier = f
// }

func (this *config_mgr) SetWatcher(w watcher.IWatcher) error {
	if this.watcher != nil {
		if err := this.watcher.Stop(); err != nil {
			fmt.Println("err:", err)
		}
	}
	this.watcher = w
	return this.watcher.SetDistributeFunc(this)
}

func (this *config_mgr) Regist(item item.IItem) error {
	if this.watcher == nil {
		return errors.New("watcher is nil")
	}
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
	if err := this.watcher.Regist(item); err != nil {
		return err
	}
	if this.checker != nil {
		// i := file_identifier
		// //设置替换
		// if this.gen_checker_identifier != nil {
		// 	i = this.gen_checker_identifier(item)
		// }
		// // 每个item都有自己的检测标识,设置自己的替换
		// if item.Opts() != nil &&
		// 	item.Opts().CheckerIdentifier != nil &&
		// 	*(item.Opts().CheckerIdentifier) != "" {
		// 	i = *(item.Opts().CheckerIdentifier)
		// }
		if err := this.checker.Regist(item); err != nil {
			return err
		}
	}
	return nil
}

func (this *config_mgr) GetLoadItem(rt reflect.Type, flag string) *item.ItemValue {
	for _, item := range this.items {
		if item.RT() == rt {
			for _, item_meta := range item.RVS() {
				if flag == item_meta.Flag {
					return item_meta
				}
			}
		}
	}
	return nil
}

// 数据会从这里分发出来
func (this *config_mgr) Distribute(file_identifier string, buf []byte) error {
	var err error
	for _, item := range this.items {
		if b, _ := item.Match(file_identifier); b {
			if err1 := item.LoadFile(file_identifier, buf); err1 != nil {
				err = err1
				continue
			}
		}
	}
	return err
}
