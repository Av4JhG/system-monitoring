package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// MaxSeconds - дефолтное время хранения метрик.
const MaxSeconds = 600

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
	LoadDisks         bool
	UsedFS            bool
	TopTalkersNetwork bool
	NetworkStat       bool
}

func NewConfig(configFile string) (Config, error) {
	var config Config

	v := viper.New()

	configure(v)
	if len(configFile) > 0 {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return config, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return config, nil
}

func configure(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.maxSeconds", MaxSeconds)
	v.SetDefault("log.level", "INFO")
	v.SetDefault("log.file", "app.log")
	v.SetDefault("server.port", "8100")
	v.SetDefault("server.host", "localhost")
	v.SetDefault("metric.loadAvg", true)
	v.SetDefault("metric.cpu", true)
	v.SetDefault("metric.loadDisks", true)
	v.SetDefault("metric.usedFs", true)
}

func (c Config) Validate() error {
	if err := c.App.Validate(); err != nil {
		return err
	}
	if err := c.Server.Validate(); err != nil {
		return err
	}

	return nil
}

func (c AppConf) Validate() error {
	if c.MaxSeconds <= 0 {
		return errors.New("time to keep metrics must be greater than zero")
	}

	return nil
}

func (c ServerConf) Validate() error {
	if c.Port == "" {
		return errors.New("server port is required")
	}

	return nil
}
