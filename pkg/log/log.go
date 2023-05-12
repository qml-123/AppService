package log

import (
	"github.com/qml-123/app_log/logger"
)

func InitLogger(url []string) (err error) {
	err = logger.NewLogger(url, "log")
	if err != nil {
		return
	}
	return nil
}
