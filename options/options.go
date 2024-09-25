package options

// 之所以全部用指针,是因为不确定是否零值是否需要被使用
type Option struct {
	Len               *int
	Cap               *int
	CheckerIdentifier *string //检测的标识可能需要随机兼容,所以自定义,不存在的话使用file_identifier

	//下面2个方法有且只有一个函数生效
	CheckerFunc      *func([]byte) error           //注入一个检查器,优先使用这个,次要使用默认load函数
	CheckerMatchFunc *func([]string, []byte) error //注入一个具有匹配参数检查器,优先使用这个,次要使用默认load函数
}

func New() *Option {
	return &Option{}
}

func (this *Option) SetCheckerIdentifier(s string) *Option {
	this.CheckerIdentifier = &s
	return this
}

func (this *Option) SetLen(i int) *Option {
	this.Len = &i
	return this
}

func (this *Option) SetCap(i int) *Option {
	this.Cap = &i
	return this
}

func (this *Option) SetCheckerFunc(f func([]byte) error) *Option {
	if f == nil {
		return this
	}
	this.CheckerFunc = &f
	return this
}

// 一般正则匹配的时候需要
func (this *Option) SetCheckerMatchFunc(f func([]string, []byte) error) *Option {
	if f == nil {
		return this
	}
	this.CheckerMatchFunc = &f
	return this
}

func (this *Option) Merge(opts ...*Option) *Option {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

func (this *Option) merge(opt *Option) {
	if opt.Len != nil {
		this.Len = opt.Len
	}
	if opt.Cap != nil {
		this.Cap = opt.Cap
	}
	if opt.CheckerIdentifier != nil {
		this.CheckerIdentifier = opt.CheckerIdentifier
	}

	if opt.CheckerFunc != nil {
		this.CheckerFunc = opt.CheckerFunc
	}

	if opt.CheckerMatchFunc != nil {
		this.CheckerMatchFunc = opt.CheckerMatchFunc
	}

}
