package watcher

import "github.com/ndsky1003/config/v2/item"

type IWatcher interface {
	Stop() error
	SetDistributeFunc(IDistributer) error
	Regist(item item.IItem) error
}

type IDistributer interface {
	// @Distribute 内容分发
	// @file_identifier 文件的唯一路径
	// @buf 文件内容
	Distribute(file_identifier string, buf []byte) error
}
