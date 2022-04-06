package zdpgo_env

import (
	"github.com/zhangdapeng520/zdpgo_env/libs/godotenv"
	"os"
)

// Env 操作环境变量的核心对象
type Env struct {
	Load func(filenames ...string) (err error)
}

func New() *Env {
	e := Env{}

	// 加载方法
	e.Load = godotenv.Load
	return &e
}

// Get 根据键获取环境变量中的值
func (e *Env) Get(key string) string {
	result := os.Getenv(key)
	return result
}

// GetDefault 获取环境变量信息，找不到则使用默认值
func (e *Env) GetDefault(key, defaultValue string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defaultValue
	}
	return val
}
