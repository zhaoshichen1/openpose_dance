package util

import (
	"fmt"
	"io/ioutil"
)

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
		fmt.Println("写入文件成功:", content)
	}
}
