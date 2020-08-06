package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

var config_file_name = "000_config.toml"

var config SConfig

func (conf SConfig) Dump() string {
	b, err := json.Marshal(conf)
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	return out.String()
}

func CreateConfigFile() {
	ctoml := SConfig{
		Port:          8181,
		IP:            "",
		MaxSize:       1 << 20,
		ExecPath:      "",
		SaveDir:       "upload_img/",
		Field:         "file",
		UploadFileUrl: "/upload_img_123",
	}
	var buf bytes.Buffer
	toml.NewEncoder(&buf).Encode(ctoml)
	//写入配置文件
	ioutil.WriteFile(config_file_name, buf.Bytes(), os.ModePerm)
}

func LoadingConfigSourceFile() (_tmp SConfig, err error) {
	for _, f := range []string{"./", "./../../"} {
		//_tmp := SConfig{}
		_, err = toml.DecodeFile(f+config_file_name, &_tmp)
		if err == nil {
			abs, e2 := filepath.Abs(f + config_file_name)
			if e2 != nil {
				return
			}
			_tmp.ExecPath = filepath.Dir(abs)
			config = _tmp
			err = nil
			break
		}
	}
	if _tmp.ExecPath == "" {
		err = errors.New("加载配置文件失败")
		return
	}
	_tmp.SaveDir = _tmp.ExecPath + "/" + _tmp.SaveDir
	return
}

func GetConfig() SConfig {
	return config
}
