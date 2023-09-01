package config

import (
	"testing"
)

func TestNew(t *testing.T) {
	Regist[[]int]("./config", nil)
	v := Get[[]int]()
	t.Log(*v == nil, *v)

	// c, _ := os.Getwd()
	// fmt.Println(filepath.Join(c, "../config/config/"))
	// v := filepath.Join(c, "config/")
	// fmt.Println(v)
	// fmt.Println(filepath.Join(c, "/Users/mac/go/workSpace/self-pkg/config/config"))
	// _ = New(
	// 	options.New().
	// 		SetDir("./config1").
	// 		SetDir(fmt.Sprintf("%v/%v", c, "config")).
	// 		SetDir("config/").SetDir("../config/config/"),
	// )
	//
	// select {}
}
