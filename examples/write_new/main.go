package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_env"
)

func main() {
	e := zdpgo_env.New()
	envMap := map[string]string{
		"a": "bbb",
		"b": "ccc",
		"c": "ddd",
	}

	// 写入环境变量
	err := e.WriteNew(".env1", envMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取环境变量
	fmt.Println(e.Get("a"))
	fmt.Println(e.Get("b"))
	fmt.Println(e.Get("c"))
}
