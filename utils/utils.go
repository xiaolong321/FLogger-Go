package utils

import (
	"os"
	"time"
)

// PathExists xxx
func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// Contain xxx
func Contain(array []int, item int) bool {
	for _, elem := range array {
		if elem == item {
			return true
		}
	}
	return false
}

// GetCurrTimestamp xxx
func GetCurrTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrTimeInCST xxx
func GetCurrTimeInCST() string {
	return time.Now().String()
}

// GetCurrTime xxx
func GetCurrTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// GetCurrDate xxx
func GetCurrDate() string {
	return time.Now().Format("2006-01-02")
}
