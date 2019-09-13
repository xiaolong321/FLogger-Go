package flogger

const (
	/** 日志类型  */
	// 调试信息
	DEBUG = iota
	// 普通信息
	INFO
	// 警告信息
	WARN
	// 错误信息
	ERROR
	// 严重错误信息
	FATAL
)

var (
	// 日志类型描述表
	LoggerLevelMap = map[int]string{
		0: "DEBUG",
		1: "INFO",
		2: "WARN",
		3: "ERROR",
		4: "FATAL",
	}
)
