package core

import (
	"os"
	"path/filepath"
	"strings"
)

const procPath = "/proc"

// ReadProcFile считывает указанный файл из /proc/ и возвращает его как набор строк.
func ReadProcFile(name string) ([]string, error) {
	content, err := os.ReadFile(procFileName(name))
	if err != nil {
		return nil, err
	}

	return SplitLines(string(content)), nil
}

// SplitLines разбивает полученный текст на слайс строк.
func SplitLines(content string) []string {
	return strings.Split(strings.TrimSpace(content), "\n")
}

func procFileName(name string) string {
	return filepath.Join(procPath, name)
}
