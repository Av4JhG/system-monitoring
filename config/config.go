package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// структура Config содержит конфигурацию приложения.
type Config struct {
	Logger LoggerConf
	App    AppConf
	Server ServerConf
	Metric MetricConf
}

// структура LoggerConf содержит настройки логирования.
type LoggerConf struct {
	Level string
	File  string
}

// структура AppConf содержит общие настройки приложения.
type AppConf struct {
	MaxSeconds int
}

// структура ServerConf содержит настройки gRPC сервера.
type ServerConf struct {
	Host string
	Port string
}

// структура MetricConf позволяет управлять сбором определенных метрик.
type MetricConf struct {
	CPU               bool
	LoadAvg           bool
	LoadDisk          bool
	UseFs             bool
	TopTalkersNetwork bool
	NetworkStat       bool
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	return &config, nil
}
