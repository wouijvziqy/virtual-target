package log

import (
	"fmt"
	"time"
)

func Info(str string) {
	fmt.Printf("[info] [%s] %s\n", getTime(), str)
}

func Error(str string) {
	fmt.Printf("[error] [%s] %s\n", getTime(), str)
}

func getTime() string {
	currentTime := time.Now().Format("15:04:05")
	return currentTime
}
