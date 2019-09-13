package strategy

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	cmap "github.com/summychou/FLogger-Go/concurrent-map"

	"github.com/summychou/FLogger-Go/utils"
)

// LogManager xxx
type LogManager struct {
	conf       *LogManagerConfig
	logFileMap cmap.ConcurrentHashMap
}

var gloLogManager *LogManager
var once sync.Once

// GetInstance xxx
func GetInstance() *LogManager {
	once.Do(func() {
		gloLogManager = new(LogManager)
		gloLogManager.logFileMap = cmap.New()
	})
	return gloLogManager
}

// SetConf xxx
func (manager *LogManager) SetConf(fileSize, cacheSize, duration int64) {
	manager.conf = new(LogManagerConfig)
	manager.conf.LogFileSize = fileSize
	manager.conf.LogCacheSize = cacheSize
	manager.conf.FlushDuration = duration
}

// AppendLogEntry xxx
func (manager *LogManager) AppendLogEntry(filename string, msg *bytes.Buffer) {
	var lfi *LogFileItem
	if value, ok := manager.logFileMap.Get(filename); !ok {
		if value, ok = manager.logFileMap.Get(filename); !ok { // 双重检测锁
			// TODO: 需要加锁吗？
			lfi = NewLogFileItem()
			lfi.logFileName = filename
			lfi.nextWriteTime = utils.GetCurrTimestamp() + manager.conf.FlushDuration
			manager.logFileMap.Set(filename, lfi)
		} else {
			lfi = value.(*LogFileItem)
		}
	} else {
		lfi = value.(*LogFileItem)
	}
	lfi.Lock()
	if lfi.currLogBuff == 'A' {
		lfi.logBufferA = append(lfi.logBufferA, msg)
	} else if lfi.currLogBuff == 'B' {
		lfi.logBufferB = append(lfi.logBufferB, msg)
	}
	lfi.currCacheSize += int64(len(msg.String()))
	lfi.Unlock()
	manager.logFileMap.Set(filename, lfi)
}

func (manager *LogManager) flush(force bool) {
	currTime := utils.GetCurrTimestamp()
	for kv := range manager.logFileMap.IterBuffered() {
		lfi := kv.Val.(*LogFileItem)
		if currTime >= lfi.nextWriteTime ||
			manager.conf.LogCacheSize <= lfi.currCacheSize ||
			force {
			var wrtLog []*bytes.Buffer
			lfi.Lock()
			if lfi.currLogBuff == 'A' {
				wrtLog = lfi.logBufferA
				lfi.logBufferA = make([]*bytes.Buffer, 0)
				lfi.currLogBuff = 'B'
			} else if lfi.currLogBuff == 'B' {
				wrtLog = lfi.logBufferA
				lfi.logBufferB = make([]*bytes.Buffer, 0)
				lfi.currLogBuff = 'A'
			}
			lfi.currLogSize += manager.writeToDisk(lfi, wrtLog)
			lfi.Unlock()
			manager.logFileMap.Set(kv.Key, lfi)
		}
	}
}

func (manager *LogManager) writeToDisk(lfi *LogFileItem, wrtLog []*bytes.Buffer) int64 {
	manager.createLogFile(lfi)
	fd, err := os.OpenFile(lfi.fullLogFileName, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	var size int64
	for _, logMsg := range wrtLog {
		if n, err := fd.Write(logMsg.Bytes()); err != nil {
			size += int64(n)
		}
	}
	return size
}

func (manager *LogManager) createLogFile(lfi *LogFileItem) {
	// 获取当前系统日期
	currDate := utils.GetCurrDate()
	// 判断日志根路径是否存在，不存在则先创建
	if !utils.PathExists(manager.conf.LogPath) {
		if err := os.MkdirAll(manager.conf.LogPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
	// 如果超过单个文件大小，则拆分文件
	if lfi.fullLogFileName != "" && lfi.currLogSize >= manager.conf.LogFileSize {
		if utils.PathExists(lfi.fullLogFileName) {
			newFileName := fmt.Sprintf("%s/%s/%s_%s.log",
				manager.conf.LogPath,
				lfi.lastWriteDate,
				lfi.logFileName,
				utils.GetCurrTime(),
			)
			if err := os.Rename(lfi.fullLogFileName, newFileName); err != nil {
				fmt.Printf("日志已成功备份为%s\n!", newFileName)
			} else {
				fmt.Println("日志备份失败!")
			}
			lfi.fullLogFileName = ""
			lfi.currLogSize = 0
		}
	}
	// 创建文件
	if lfi.fullLogFileName == "" || lfi.lastWriteDate != currDate {
		path := fmt.Sprintf("%s/%s", manager.conf.LogPath, currDate)
		if !utils.PathExists(path) {
			if err := os.Mkdir(manager.conf.LogPath, os.ModePerm); err != nil {
				panic(err)
			}
		}
		lfi.fullLogFileName = fmt.Sprintf("%s/%s.log", path, lfi.logFileName)
		lfi.lastWriteDate = currDate
		if utils.PathExists(lfi.fullLogFileName) {
			info, _ := os.Stat(lfi.fullLogFileName)
			lfi.currLogSize = info.Size() // TODO: ?
		} else {
			lfi.currLogSize = 0
		}
	}
}
