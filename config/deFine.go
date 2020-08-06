package config

//定义 系统基础设置
type SConfig struct {

	//端口
	Port int

	//监听IP
	IP string

	//接口地址
	UploadFileUrl string

	//允许文件最大值
	MaxSize int64

	ExecPath string

	//保存位置
	SaveDir string

	//上传文件的字段名字
	Field string
}
