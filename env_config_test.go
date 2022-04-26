package zdpgo_env

import (
	"github.com/zhangdapeng520/zdpgo_test"
	"testing"
)

type envConfig struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}

// 测试解析配置
func TestEnv_ParseConfig(t *testing.T) {
	var cfg envConfig
	e := getEnv()
	e.Read(".env")      // 读取配置，写入到环境变量
	e.ParseConfig(&cfg) // 解析配置为结构体

	// 创建断言对象
	test := zdpgo_test.NewWirthConfig(zdpgo_test.Config{TestObj: t})

	// 创建表格
	var tests = []struct {
		host string // 期望输出
		port int
	}{
		{"127.0.0.1", 8888},
	}

	// 遍历表格数据
	for _, testData := range tests {
		// 断言
		test.Assert.Equal(cfg.Host, testData.host)
		test.Assert.Equal(cfg.Port, testData.port)
	}
}
