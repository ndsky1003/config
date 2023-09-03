package options

type Option struct {
	Len *int
	Cap *int
}

func New() *Option {
	return &Option{}
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
}
