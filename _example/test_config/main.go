package main

import (
	"fmt"
	"github.com/bbs-v2ex/img_server/config"
)

func main() {
	config.CreateConfigFile()
	file, err := config.LoadingConfigSourceFile()
	fmt.Println(err, file.Dump())
}
