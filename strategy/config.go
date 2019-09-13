package strategy

// LogManagerConfig xxx
type LogManagerConfig struct {
	LogPath       string `yaml:"log_path"`
	LogFileSize   int64  `yaml:"log_file_size"`
	LogCacheSize  int64  `yaml:"log_cache_size"`
	FlushDuration int64  `yaml:"flush_duration"`
}
