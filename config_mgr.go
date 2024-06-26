package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ndsky1003/config/checker"
	"github.com/ndsky1003/config/item"
	"github.com/ndsky1003/config/watcher"
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
	items                   []item.IItem
	watcher                 watcher.IWatcher
	checker                 checker.IChecker
	checker_identifier_func func(string) string
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

func (this *config_mgr) SetCheckerIdentifierFunc(f func(string) string) {
	this.checker_identifier_func = f
}

func (this *config_mgr) SetWatcher(w watcher.IWatcher) {
	if this.watcher != nil {
		if err := this.watcher.Stop(); err != nil {
			fmt.Println("err:", err)
		}
	}
	this.watcher = w
	_ = this.watcher.SetDistributeFunc(this.DistributeData)
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
		i := file_identifier
		//设置替换
		if this.checker_identifier_func != nil {
			i = this.checker_identifier_func(i)
		}
		// 每个item都有自己的检测标识,设置自己的替换
		if item.Opts() != nil &&
			item.Opts().CheckerIdentifier != nil &&
			*(item.Opts().CheckerIdentifier) != "" {
			i = *(item.Opts().CheckerIdentifier)
		}
		this.checker.On(i, item.CheckBuf)
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

// MARK check 也可注入该方法检查，成功才会替换
func (this *config_mgr) DistributeData(file_identifier string, buf []byte) error {
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
