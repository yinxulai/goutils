package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// ReadAll 读取全部
func ReadAll(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return ioutil.ReadAll(file)
}

// ReadBlock 读取指定块大小
func ReadBlock(filePath string, bufSize int, callback func([]byte) bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	buf := make([]byte, bufSize) // 一次读取多少个字节
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
func ReadLine(filePath string, callback func([]byte) bool) error {
	file, err := os.Open(filePath)
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
func WriteByte(filePath string, append bool, data []byte) error {
	var err error
	var file *os.File
	defer file.Close()

	// 文件是否存在
	exits, err := PathExists(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !exits { // 如果文件不存在 创建
		// 创建目录
		dir := path.Dir(filePath)
		if dir != "." || dir != "/" {
			err = os.MkdirAll(dir, 0666)
			if err != nil {
				return err
			}
		}

		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if exits && append { // 如果文件存在 且为追加 则追加模式打开
		file, err = os.OpenFile(filePath, os.O_APPEND, 0666)
		if err != nil {
			return err
		}
	}

	if exits && !append { // 如果文件存在 读写模式打开
		file, err = os.OpenFile(filePath, os.O_WRONLY, 0666)
		if err != nil {
			return err
		}

		// 清空文件先
		err = file.Truncate(0)
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
func WriteString(filePath string, append bool, data string) error {
	return WriteByte(filePath, append, []byte(data))
}

// PathExists 文件是否存在
func PathExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
