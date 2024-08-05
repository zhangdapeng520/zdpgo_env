package zdpgo_env

import "github.com/zhangdapeng520/zdpgo_env/env"

func EnvToStruct(v interface{}) error {
	return env.Parse(v)
}
