package checker

// file_identifier1 文件标示符，可以是正则，无规则限制 ,Regist的时候的值
// buf 文件内容
// err 返回值，当存在返回值的时候，自行处理
type IChecker interface {
	On(file_identifier1 string, fn func(buf []byte) error)
}
