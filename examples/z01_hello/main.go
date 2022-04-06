package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_env"
	"log"
)

func main() {
	e := zdpgo_env.New()
	err := e.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	s3Bucket := e.Get("S3_BUCKET")
	secretKey := e.Get("SECRET_KEY")

	fmt.Println(s3Bucket, secretKey)

	// 读取系统环境变量
	fmt.Println(e.Get("GOPATH"))
}
