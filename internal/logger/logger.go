package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/Av4JhG/system-monitoring/internal/sm"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(logLevel string) (sm.Logger, error) {
	// Создаём папку logs, если не существует.
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0o755)
	}

	// Настраиваем lumberjack для ротации логов.
	log := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50,   // Мб.
		MaxBackups: 30,   // Кол-во файлов.
		MaxAge:     30,   // Дней.
		Compress:   true, // Сжимать старые.
	}
	level, err := parseLevel(logLevel)
	if err != nil {
		level = slog.LevelInfo
	}

	// Изменяем формат веремени в логе.
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}

			t := a.Value.Time()

			a.Value = slog.StringValue(t.Format(time.DateTime))

			return a
		},
	}
	// Обёртка логгера
	jsonLogger := slog.New(slog.NewJSONHandler(log, opts))

	resultLogger := logger{
		logger: jsonLogger,
	}

	return resultLogger, nil
}

func parseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
}

type logger struct {
	logger *slog.Logger
}

// Debug логирует сообщения отладки.
func (l logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info логирует информационные сообщения.
func (l logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Error логирует ошибки.
func (l logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}
