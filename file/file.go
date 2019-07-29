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

// WriteByte 写入文件
// TODO: 自动创建文件时，目前不支持自动创建目录，例如 path 为 ./dir/data.json
// 无法自动创建 dir 目录
func WriteByte(path string, append bool, data []byte) error {
	var err error
	var file *os.File
	defer file.Close()

	exits, err := PathExists(path)
	if err != nil {
		return err
	}

	if !exits { // 如果文件不存在 创建
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	}

	if exits && !append { // 如果文件存在 读写模式打开
		file, err = os.OpenFile(path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	}

	if exits && append { // 如果文件存在 且为追加 则追加模式打开
		file, err = os.OpenFile(path, os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// WriteString 写入字符串
func WriteString(path string, append bool, data string) error {
	return WriteByte(path, append, []byte(data))
}

// PathExists 文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
