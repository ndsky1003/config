# config

RCU机制的注入与获取

#### 原则就是一个文件对应一个类型

一个reflect.Type
一个文件
以上2个均可找到对应的回调方法

usage
```go

type Person struct {
	Name string `json:"Name"`
}

func load3(buf []byte) (*[]*Person, error) {
	var v []*Person
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

if err := Regist[[]*Person]("./config9/reg:person_([a-z]{3})_\\d*.yaml", load3); err != nil {
	panic(err)
}

v4 := Get[[]*Person]("bbo") 
v4 is data

```

#### TODO
找个合适的日志库
