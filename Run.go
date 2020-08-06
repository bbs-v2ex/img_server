package img_server

import (
	"encoding/json"
	"fmt"
	"github.com/bbs-v2ex/img_server/config"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
)

const (
	PNG  = "png"
	JPG  = "jpg"
	JPEG = "jpeg"
	GIF  = "gif"
)

var _con = config.SConfig{}
var err error

type result struct {
	Url     string `json:"url"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Size    int64  `json:"size"`
}

func Server() {
	_con, err = config.LoadingConfigSourceFile()
	if err != nil {
		log.Fatal("读取配置文件失败1111", err)
	}

	//判断文件夹是否存在如果不存在则创建
	dir, _ := IsDir(_con.SaveDir)

	if !dir {
		err2 := os.MkdirAll(_con.SaveDir, os.ModePerm)
		if err2 != nil {
			log.Fatal("读取配置文件失败3333", err)
		}
	}

	http.HandleFunc("/upload", upload)

	http.HandleFunc("/", outimg)

	loca_url := fmt.Sprintf("%s:%d", _con.IP, _con.Port)
	fmt.Println(loca_url)
	err := http.ListenAndServe(loca_url, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func to(i interface{}) []byte {
	marshal, _ := json.Marshal(i)
	return marshal
}
