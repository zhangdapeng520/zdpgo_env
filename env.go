package zdpgo_env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Env 操作环境变量的核心对象
type Env struct {
	Config       *Config                      // 核心配置
	FileEnvMap   map[string]map[string]string // 文件环境变量
	SystemEnvMap map[string]string            // 系统环境变量
	AllEnvMap    map[string]string            // 所有环境变量
}

// New 创建环境变量实例
func New() (e *Env) {
	cfg := Config{}
	e = NewWithConfig(cfg)
	return
}

// NewWithConfig 使用配置创建环境变量实例
func NewWithConfig(config Config) (e *Env) {
	e = &Env{}
	e.Config = &config
	e.SystemEnvMap = make(map[string]string)
	e.AllEnvMap = make(map[string]string)
	e.FileEnvMap = make(map[string]map[string]string)
	return e
}

// Parse 解析带env标签的结构体，并将其值添加到环境变量中
// @param v 结构体对象
// @param opts 参数选项
// @return opts 参数选项
func (e *Env) Parse(v interface{}, opts ...Options) (err error) {
	err = ParseConfig(v, opts...)
	return
}

// ReadFiles 加载配置文件列表中的所有环境变量
// @param filenames 环境变量文件列表
func (e *Env) ReadFiles(filenames ...string) (err error) {
	err = e.Load(filenames...)
	e.FindAll() // 将所有环境变量同步到EnvMap中
	return
}

// ReadDefault 读取默认的环境变量配置文件
func (e *Env) ReadDefault() (err error) {
	err = e.Read(".env")
	e.FindAll() // 将所有环境变量同步到EnvMap中
	return
}

// Read 读取环境变量配置文件
// @param filename 环境变量文件名称
// @return err 错误信息
func (e *Env) Read(filename string) (err error) {
	return e.ReadFiles(filename)
}

// Save 保存环境变量
func (e *Env) Save(filename string) (err error) {
	if envMap, ok := e.FileEnvMap[filename]; ok {
		return Write(envMap, filename)
	}
	err = errors.New(fmt.Sprintf("不存在该文件`%s`的配置\n", filename))
	return
}

// Write 保存环境变量
func (e *Env) Write(filename string) (err error) {
	return e.Save(filename)
}

// WriteNew 将环境变量map写入到新的配置文件中
func (e *Env) WriteNew(filename string, envMap map[string]string) (err error) {
	// 将数据写入到环境变量的配置文件
	err = Write(envMap, filename)
	if err != nil {
		return
	}

	// 重新将配置文件中的环境变量读取并写入系统环境变量
	err = e.Read(filename)
	if err != nil {
		return
	}
	return
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

// Add 添加环境变量
func (e *Env) Add(key, value string) (err error) {
	e.AllEnvMap[key] = value
	err = os.Setenv(key, value)
	return
}

// AddMany 同时添加多个环境变量
func (e *Env) AddMany(keyValues ...string) (err error) {
	if keyValues == nil || len(keyValues) == 0 {
		err = errors.New("要添加的环境变量不能为空")
		return
	}
	if len(keyValues)%2 != 0 {
		err = errors.New("要添加的环境变量键值对不匹配")
		return
	}

	for i := 0; i < len(keyValues); i += 2 {
		err = e.Add(keyValues[i], keyValues[i+1])
		if err != nil {
			return
		}
	}
	return
}

// Find 查找环境变量
func (e *Env) Find(key string) (value string, exists bool) {
	value, exists = os.LookupEnv(key)
	return
}

// Remove 移除环境变量
func (e *Env) Remove(key string) (err error) {
	delete(e.AllEnvMap, key)
	err = os.Unsetenv(key)
	return
}

// RemoveMany 同时移除多个环境变量
func (e *Env) RemoveMany(keys ...string) (err error) {
	if keys == nil || len(keys) == 0 {
		err = errors.New("要移除的环境变量不能为空")
		return
	}

	for _, key := range keys {
		err = e.Remove(key)
		if err != nil {
			return
		}
	}
	return
}

// FindAll 读取所有的环境变量
func (e *Env) FindAll() (result map[string]string) {
	result = make(map[string]string)
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		key := envPair[0]
		value := envPair[1]
		result[key] = value
		e.AllEnvMap[key] = value
	}
	return
}
