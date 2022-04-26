package zdpgo_env

import (
	"fmt"
	"testing"
)

func getEnv() *Env {
	return New()
}

func TestEnv_basic(t *testing.T) {
	e := getEnv()

	// 加载.env中的环境变量
	e.ReadDefault()

	// 获取环境变量
	host := e.Get("host")
	fmt.Println(host)

	// 设置环境变量
	e.Add("password", "abc123")
	fmt.Println(e.Get("password"))

	// 移除环境变量
	e.Remove("password")
	fmt.Println(e.Get("password"))
}

func TestEnv_FindAll(t *testing.T) {
	e := getEnv()

	// 加载.env中的环境变量
	e.Read(".env")

	// 获取环境变量
	keyValues := e.FindAll()
	for k, v := range keyValues {
		fmt.Println(k, v)
	}
}

func TestEnv_ReadWrite(t *testing.T) {
	e := getEnv()

	// 加载.env中的环境变量
	e.ReadDefault()

	// 修改
	e.FileEnvMap[".env"]["host"] = "127.0.0.1"

	// 写入环境变量
	err := e.Write(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestEnv_WriteNew(t *testing.T) {
	e := getEnv()

	envMap := map[string]string{
		"a": "bbb",
		"b": "ccc",
	}

	// 写入环境变量
	err := e.WriteNew(".env1", envMap)
	if err != nil {
		fmt.Println(err)
		return
	}
}
