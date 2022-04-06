package zdpgo_env

import (
	"github.com/zhangdapeng520/zdpgo_env/libs/envconfig"
	"github.com/zhangdapeng520/zdpgo_env/libs/godotenv"
	"os"
)

// Env 操作环境变量的核心对象
type Env struct {
	// 加载配置文件中的变量到环境变量
	Load func(filenames ...string) (err error)

	// 读取环境变量到配置对象
	Read func(config interface{}, opts ...envconfig.Options) error
}

func New() *Env {
	e := Env{}

	// 加载方法
	e.Load = godotenv.Load

	// 读取配置的方法
	e.Read = envconfig.Parse

	// 返回环境变量对象
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
