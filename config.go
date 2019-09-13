package flogger

import (
	"github.com/summychou/FLogger-Go/strategy"
)

// FLoggerConfig xxx
type FLoggerConfig struct {
	LoggerLevel  []int `yaml:"logger_level"`
	ConsolePrint bool  `yaml:"console_print"`

	strategy.LogManagerConfig
}
