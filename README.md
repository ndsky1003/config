# config

RCU机制的注入与获取

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

	v4 := config.Get[[]*Person]()
	fmt.Printf("%+v", v4)
}

```

#### 2.正则使用
1. 文件名以reg:开头
2. 正则命里需要用（）来提取标识符
3. 标识符用`_`连接,作为最终获取的
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
	if err := config.Regist[[]*Person]("./config9/reg:person_([a-z0-9]{0,10})_(\\d*).yaml", load3); err != nil {
		panic(err)
	}

	v4 := config.Get[[]*Person]("hilo2l_2")
	fmt.Printf("%+v", v4)
}

```


```go
//包装现有的cfgmgr
	if err := config.Regist[ConfigVip]("reg:db_vip_([a-z]{3})_([a-z]{3}).yaml", loadConfigVip); err != nil {
		logger.Err(err)
	}

```
