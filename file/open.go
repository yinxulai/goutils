package file

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

// ReadAll 读取全部
func ReadAll(filePth string) ([]byte, error) {
	file, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(file)
}

// ReadBlock 读取指定块大小
func ReadBlock(filePth string, bufSize int, callback func([]byte) bool) error {
	file, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, bufSize) //一次读取多少个字节
	bfRd := bufio.NewReader(file)
	for {
		n, err := bfRd.Read(buf)
		if err != nil {
			//遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
		// n 是成功读取字节数
		if !callback(buf[:n]) {
			// 返回 false 停止读取
			break
		}
	}
	return nil
}

// ReadLine 读取行
func ReadLine(filePth string, callback func([]byte) bool) error {
	file, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer file.Close()

	bfRd := bufio.NewReader(file)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}

		if !callback(line) {
			// 返回 false 停止读取
			break
		}
	}
	return nil
}
