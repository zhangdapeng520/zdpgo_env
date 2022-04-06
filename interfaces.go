package zdpgo_env

import "github.com/zhangdapeng520/zdpgo_env/libs/envconfig"

// Load 加载文件接口
type Load interface {
	// Load 加载文件中定义的环境变量到系统环境变量中
	Load(filenames ...string) (err error)
}

type Reader interface {
	// Read 读取环境变量进配置中
	Read(config interface{}, opts ...envconfig.Options) error
}
