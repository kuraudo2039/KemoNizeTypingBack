package util

import (
	"fmt"
	"os"
)

type LogObj struct {
	Message string
	Data    interface{}
}

// 環境変数 ENVIRONMENT に従って、ログの出力を制御する
func Log(log LogObj) {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "DEV":
		fmt.Println(log.Message)
		if log.Data != nil {
			fmt.Println(log.Data)
		}
	case "PROD":
	default:
	}
}
