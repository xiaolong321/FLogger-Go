package strategy

import (
	"bytes"
	"sync"
)

// LogFileItem xxx
type LogFileItem struct {
	logFileName     string          // 不包括路径且不带扩展名的日志文件名称
	fullLogFileName string          // 包括路径且带扩展名的日志文件名称
	currLogSize     int64           // 当前日志文件大小
	currCacheSize   int64           // 当前已缓存大小
	currLogBuff     byte            // 当前正在使用的日志缓冲队列
	logBufferA      []*bytes.Buffer // 日志缓冲队列A
	logBufferB      []*bytes.Buffer // 日志缓冲队列B
	lastWriteDate   string          // 上次写入时的日期
	nextWriteTime   int64           // 下次日志输出到文件时间戳

	sync.Mutex
}

// NewLogFileItem xxx
func NewLogFileItem() *LogFileItem {
	return &LogFileItem{
		currLogBuff: 'A',
		logBufferA:  make([]*bytes.Buffer, 0),
		logBufferB:  make([]*bytes.Buffer, 0),
	}
}
