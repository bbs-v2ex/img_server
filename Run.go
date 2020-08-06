package img_server

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/bbs-v2ex/img_server/config"
	"image"
	"image/gif"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
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

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	res := result{}
	if r.Method != "POST" {
		res.Message = "请求类型需要是 POST"
		res.Code = 10001
		w.Write(to(res))
		return
	}

	//缓冲的大小 - 4M

	r.ParseMultipartForm(1024 << 12)
	//数据大小的限制
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	upfile, upFileInfo, err := r.FormFile(_con.Field)
	if err != nil {
		res.Message = "提取上传的文件失败"
		res.Code = 10002
		w.Write(to(res))
		return
	}
	defer upfile.Close()
	bufUpFile := bufio.NewReader(upfile)
	//进行图片的解码
	img, imgtype, err := image.Decode(bufUpFile)
	if err != nil {
		res.Message = err.Error()
		res.Code = 10002
		w.Write(to(res))
		return
	}
	_, err = upfile.Seek(0, 0)
	if err != nil {
		res.Message = err.Error()
		res.Code = 10003
		w.Write(to(res))
		return
	}

	//计算文件 md5 值

	md5Hash := md5.New()
	// 读入缓存
	bufFile := bufio.NewReader(upfile)
	_, err = io.Copy(md5Hash, bufFile)
	if err != nil {
		res.Message = err.Error()
		res.Code = 10004
		w.Write(to(res))
		return
	}
	fileMd5FX := md5Hash.Sum(nil)
	fileMd5 := fmt.Sprintf("%x", fileMd5FX)
	filePath := _con.SaveDir + fileMd5
	dirPath := filePath
	filePath += "/" + fileMd5

	// 获取目录信息，并创建目录
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		err = os.MkdirAll(dirPath, 0666)
		if err != nil {

			res.Message = err.Error()
			res.Code = 10005
			w.Write(to(res))
			return
		}
	} else {
		if !dirInfo.IsDir() {
			err = os.MkdirAll(dirPath, 0666)
			if err != nil {
				res.Message = err.Error()
				res.Code = 10006
				w.Write(to(res))
				return
			}
		}
	}

	//
	//// 存入文件 --------------------------------------
	//
	_, err = os.Stat(filePath)
	if err != nil {
		// 打开一个文件,文件不存在就会创建
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			res.Message = err.Error()
			res.Code = 10007
			w.Write(to(res))

			return
		}
		defer file.Close()

		switch imgtype {
		case PNG:
			err = png.Encode(file, img)
			break
		case JPG, JPEG:
			err = jpeg.Encode(file, img, nil)
			break
		case GIF:

			// 重新对 gif 格式进行解码
			// image.Decode 只能读取 gif 的第一帧
			// 设置下次读写位置（移动文件指针位置）
			_, err = upfile.Seek(0, 0)
			if err != nil {
				res.Message = err.Error()
				res.Code = 10008
				w.Write(to(res))
				return
			}

			gifimg, giferr := gif.DecodeAll(upfile)
			if giferr != nil {

				res.Message = giferr.Error()
				res.Code = 10009
				w.Write(to(res))

				return
			}
			err = gif.EncodeAll(file, gifimg)
		}

		if err != nil {
			res.Message = err.Error()
			res.Code = 10010
			w.Write(to(res))
			return
		}
		res.Message = "OK"
		res.Url = fileMd5
		res.Code = 1
		res.Size = upFileInfo.Size
		w.Write(to(res))
		return
	}

	res.Message = "OK"
	res.Url = fileMd5
	res.Code = 1
	res.Size = upFileInfo.Size
	w.Write(to(res))
}
