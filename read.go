package zdpgo_env

import (
	"os"
	"strings"
)

// Load 加载.env文件中的环境变量
func (e *Env) Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = e.loadFile(filename, false)
		if err != nil {
			return
		}
	}
	return
}

// 加载配置文件中的环境变量
func (e *Env) loadFile(filename string, overload bool) (err error) {
	var (
		envMap = make(map[string]string)
	)

	// 创建文件对应的环境变量map
	e.FileEnvMap[filename] = make(map[string]string)

	// 读取文件中配置，转换为map
	envMap, err = readFile(filename)
	if err != nil {
		return
	}

	// 当前的环境变量
	currentEnv := map[string]bool{}

	// 系统的环境变量
	rawEnv := os.Environ()

	// 将系统当前的环境变量存储到map中
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	// 遍历读取到的文件中的环境变量
	for key, value := range envMap {
		e.FileEnvMap[filename][key] = value // 将数据存储进文件对应的map

		// 如果当前环境变量中不存在该key，或者是覆盖写
		if !currentEnv[key] || overload {
			os.Setenv(key, value) // 写入环境变量
		}
	}

	return nil
}
