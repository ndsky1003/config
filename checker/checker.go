package checker

import "github.com/ndsky1003/config/v2/item"

// item 文件标示符，可以是正则，无规则限制 ,Regist的时候的值
// err 返回值，当存在返回值的时候，自行处理
type IChecker interface {
	Regist(item.IItem) error
}
