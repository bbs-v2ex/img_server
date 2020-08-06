### 配置文件 `000_config.toml`
```go
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
```

### 上传测试
`curl --location --request POST 'http://127.0.0.1:8181/upload_img_123' \
 --form 'file=@/C:/Users/Administrator/Desktop/QQ截图20200804005530.png'`
 
### 显示
```go
	_img_w := r.URL.Query().Get("w")
	_w, _ := strconv.Atoi(_img_w)
	_img_h := r.URL.Query().Get("h")
	_h, _ := strconv.Atoi(_img_h)
	file_url_name := _con.SaveDir + "/" + r.URL.Path + "/" + r.URL.Path
	//如果是原文件则直接输出
```

> url 后面直接增加 w 和 h 参数 为0 则 选择图片本身的大小
 

 例如 `http://127.0.0.1:8181/b696489e8d2f095e140b682df1cd7eac?w=2000&h=100`