package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	dir, err := ioutil.ReadDir("upload_img/_avatar")
	if err != nil {
		return
	}
	l_v := []string{}
	for _, v := range dir {
		f_name := v.Name()
		if !strings.Contains(f_name, "-w") {
			l_v = append(l_v, "/_avatar/"+f_name)
		}
	}
	marshal, err := json.Marshal(l_v)
	fmt.Println(string(marshal))
}
