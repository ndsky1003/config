# config

RCU机制的注入与获取

```bash
go get github.com/ndsky1003/config/v2
```

#### 原则就是一个文件对应一个类型

一个reflect.Type
一个文件
以上2个均可找到对应的回调方法

usage
#### 1.普通使用
person.yaml
```yaml
- { Name: "d" }
- { Name: "d2" }
```

main.go
```go

type Person struct {
	Name string `yaml:"Name"`
}

func (this *Person) String() string {
	return fmt.Sprintf("%+v", *this)
}

func load3(buf []byte) (*[]*Person, error) {
	var v []*Person
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func main() {
	if err := config.Regist[[]*Person]("./config9/person.yaml", load3); err != nil {
		panic(err)
	}

	config.Get[[]*Person]()
    // v4 是一个指针
    // 有个语法糖直接将引用类型的指针去掉了
    config.GetRef([]*Person)()
}

```

#### 2.正则使用
1. 正则命里需要用（）来提取标识符
2. 标识符用`_`连接,作为最终获取的
```go
type Person struct {
	Name string `yaml:"Name"`
}

func (this *Person) String() string {
	return fmt.Sprintf("%+v", *this)
}

func load3(math []string,buf []byte) (*[]*Person, error) {
	var v []*Person
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func main() {
	if err := config.RegistByRegFunc[[]*Person]("./config9/person_([a-z0-9]{0,10})_(\\d*).yaml", load3); err != nil {
		panic(err)
	}

	v4 := config.Get[[]*Person]("hilo2l_2")
	fmt.Printf("%+v", v4)
}

```
###TODO
1. 需要Regist注册覆盖之前的注册,方便自模块重写
