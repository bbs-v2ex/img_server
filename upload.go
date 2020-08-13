package img_server

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	res := result{}
	if r.Method != "POST" {
		res.Message = "请求类型需要是 POST"
		res.Code = 10001
		w.Write(to(res))
		return
	}

	//缓冲的大小 - 4M

	//r.ParseMultipartForm(1024 << 12)
	//数据大小的限制
	r.Body = http.MaxBytesReader(w, r.Body, _con.MaxSize)
	upfile, upFileInfo, err := r.FormFile(_con.Field)
	if err != nil {
		res.Message = "提取上传的文件失败" + err.Error()
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
