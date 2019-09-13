package flogger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/summychou/FLogger-Go/strategy"
	"github.com/summychou/FLogger-Go/utils"
)

// FLogger xxx
type FLogger struct {
	conf    *FLoggerConfig
	manager *strategy.LogManager
}

var gloFLogger *FLogger
var once sync.Once

// GetInstance xxx
func GetInstance() *FLogger {
	once.Do(func() {
		gloFLogger = new(FLogger)
		gloFLogger.conf = new(FLoggerConfig)
		ymlContent, err := ioutil.ReadFile("flogger.yml")
		if err != nil {
			panic(err)
		}
		if err = yaml.Unmarshal(ymlContent, gloFLogger.conf); err != nil {
			panic(err)
		}
		gloFLogger.manager = strategy.GetInstance()
		gloFLogger.manager.SetConf(
			gloFLogger.conf.LogFileSize,
			gloFLogger.conf.LogCacheSize,
			gloFLogger.conf.FlushDuration,
		)
	})
	return gloFLogger
}

func (logger *FLogger) writeLog(filename string, level int, msg string) {
	if msg != "" && utils.Contain(logger.conf.LoggerLevel, level) {
		buffer := new(bytes.Buffer)
		logMsg := fmt.Sprintf("[%s] %s [%s] %s",
			LoggerLevelMap[level],
			getCurrTime(),
			"main", // TODO: 动态获取所在协程
			msg)
		if _, err := buffer.WriteString(logMsg); err != nil {
			panic(err)
		}
		logger.manager.AppendLogEntry(filename, buffer)
		if logger.conf.ConsolePrint && (level == ERROR || level == FATAL) {
			fmt.Println(logMsg)
		}
	}
}

// Debug 写调试日志
func (logger *FLogger) Debug(msg string) {
	logger.writeLog("debug", DEBUG, msg)
}

// Info 写普通日志
func (logger *FLogger) Info(msg string) {
	logger.writeLog("info", INFO, msg)
}

// Warn 写警告日志
func (logger *FLogger) Warn(msg string) {
	logger.writeLog("warn", WARN, msg)
}

// Error 写错误日志
func (logger *FLogger) Error(msg string) {
	logger.writeLog("error", ERROR, msg)
}

// Fatal 写严重错误日志
func (logger *FLogger) Fatal(msg string) {
	logger.writeLog("fatal", FATAL, msg)
}
