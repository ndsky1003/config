package options

type Option[T any] struct {
	Len               *int
	Cap               *int
	CheckerIdentifier *string   //检测的标识可能需要随机兼容,所以自定义,不存在的话使用file_identifier
	SuccessFunc       *func(*T) //检测成功的回调,因为加载的那个回调会调用2次,一次检测,一次加载，所以需要回调
}

func New[T any]() *Option[T] {
	return &Option[T]{}
}

func (this *Option[T]) SetCheckerIdentifier(s string) *Option[T] {
	this.CheckerIdentifier = &s
	return this
}

func (this *Option[T]) SetLen(i int) *Option[T] {
	this.Len = &i
	return this
}

func (this *Option[T]) SetCap(i int) *Option[T] {
	this.Cap = &i
	return this
}

func (this *Option[T]) SetSuccessFunc(Func func(*T)) *Option[T] {
	this.SuccessFunc = &Func
	return this
}

func (this *Option[T]) Merge(opts ...*Option[T]) *Option[T] {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

func (this *Option[T]) merge(opt *Option[T]) {
	if opt.Len != nil {
		this.Len = opt.Len
	}
	if opt.Cap != nil {
		this.Cap = opt.Cap
	}
	if opt.CheckerIdentifier != nil {
		this.CheckerIdentifier = opt.CheckerIdentifier
	}

	if opt.SuccessFunc != nil {
		this.SuccessFunc = opt.SuccessFunc
	}

}
