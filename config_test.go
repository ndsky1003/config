package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	SM  map[string]any
	SM2 map[string]any
)

type Person struct {
	Name string `json:"Name"`
}

func (this *Person) String() string {
	return fmt.Sprintf("%+v", *this)
}

func TestMain(m *testing.M) {
	if err := Regist[map[string]any]("./config/cc1.json", load); err != nil {
		panic(err)
	}
	if err := Regist[SM]("./config2/cc1.json", load2); err != nil {
		panic(err)
	}

	if err := Regist[[]*Person]("./config9/reg:person_([a-z]{3})_\\d*.yaml", load3); err != nil {
		panic(err)
	}
	m.Run()
	os.Exit(0)
}

func load3(buf []byte) (*[]*Person, error) {
	var v []*Person
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func Test_Pserson(t *testing.T) {
	for {
		time.Sleep(1e9)
		v := Get[[]*Person]("hoh_1")
		fmt.Printf("vvv:%+v\n", v)
		v1 := Get[[]*Person]()
		fmt.Printf("vvv:%+v\n", v1)
		v3 := Get[[]*Person]("hoh")
		fmt.Printf("vvv:%+v\n", v3)
		v4 := Get[[]*Person]("bbo")
		fmt.Printf("vvv:%+v\n", v4)
	}
}

func Test_nomal(t *testing.T) {
	time.Sleep(2e9)
	v := Get[map[string]any]()
	t.Log("./config/cc1.json:", v)
	v1 := Get[SM]()
	t.Logf("=================:%+v\n", v1)
}

func load(buf []byte) (*map[string]any, error) {
	var a map[string]any
	if err := json.Unmarshal(buf, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func load2(buf []byte) (*SM, error) {
	var a SM
	if err := json.Unmarshal(buf, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
