package utils

import (
	"log"
	"os"
)

// 初始化日志文件
var logFile, _ = os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
var Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

// 记录信息日志
func Info(msg string, args ...any) {
	Logger.Printf("[INFO] "+msg, args...)
}

// 记录错误日志
func LogError(msg string, args ...any) {
	Logger.Printf("[ERROR] "+msg, args...)
}
