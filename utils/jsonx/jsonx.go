package jsonx

import (
	"encoding/json"
	"os"
)

// ToJsonString 把结构体转换为json string格式
func ToJsonString(v any) (string, error) {
	dataByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(dataByte), nil
}

// ReadJsonFile 从json文件中读取数据
func ReadJsonFile(filePath string, v any) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, v); err != nil {
		return err
	}
	return nil
}
