package img_server

import (
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func outimg(w http.ResponseWriter, r *http.Request) {
	img_404 := []byte("Image resource does not exist")
	if !regexp.MustCompile(`[0-9a-zA-Z]{32}`).MatchString(r.URL.Path) {
		w.WriteHeader(404)
		w.Write(img_404)
		return
	}
	_img_w := r.URL.Query().Get("w")
	_w, _ := strconv.Atoi(_img_w)
	_img_h := r.URL.Query().Get("h")
	_h, _ := strconv.Atoi(_img_h)
	file_url_name := _con.SaveDir + "/" + r.URL.Path + "/" + r.URL.Path
	//如果是原文件则直接输出
	if _w == 0 && _h == 0 {
		file, err := os.Open(file_url_name)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Image resource does not exist 11111"))
			return
		}
		defer file.Close()

		io.Copy(w, file)

		return
	}

	file_name := fmt.Sprintf("%s-w-%d-h-%d.jpg", file_url_name, _w, _h)
	//判断文件
	if IsFile(file_name) {
		file, err := os.Open(file_name)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("Image resource does not exist 22222"))
			return
		}
		defer file.Close()

		io.Copy(w, file)
		return
	}
	img, err := imaging.Open(file_url_name)
	//如果文件存在则直接输出
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Image resource does not exist 33333"))
		return
	}
	//不能超过本来的大小
	if _w > img.Bounds().Size().X && _w != 0 {
		_w = img.Bounds().Size().X
	}
	if _h > img.Bounds().Size().Y && _h != 0 {
		_h = img.Bounds().Size().Y
	}

	dstImage128 := imaging.Resize(img, _w, _h, imaging.Lanczos)

	err = imaging.Save(dstImage128, file_name)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Image resource does not exist 4444"))
		return
	}
	imaging.Clone(img)
	//输出图像
	file, err := os.Open(file_name)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("Image resource does not exist 5555"))
		return
	}
	defer file.Close()

	io.Copy(w, file)
}
