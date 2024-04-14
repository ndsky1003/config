package options

type Option struct {
	Len               *int
	Cap               *int
	CheckerIdentifier *string //检测的标识可能需要随机兼容,所以自定义,不存在的话使用file_identifier
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
}
