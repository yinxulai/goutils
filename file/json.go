package file

import (
	"encoding/json"
)

// ReadJSON 读取 JSON
func ReadJSON(filePth string, dest interface{}) error {
	var err error

	data, err := ReadAll(filePth)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON 写入 JSON
func WriteJSON(path string, append bool, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return WriteByte(path, append, jsonData)
}
