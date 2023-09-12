package watcher

import "github.com/ndsky1003/config/item"

type IWatcher interface {
	Stop() error
	SetDistributeFunc(f func(file_identifier string, buf []byte) error) error
	Regist(item item.IItem) error
}
