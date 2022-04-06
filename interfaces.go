package zdpgo_env

// Load 加载文件接口
type Load interface {
	// Load 加载文件中定义的环境变量到系统环境变量中
	Load(filenames ...string) (err error)
}
