package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// I MAKE THIS FORMATING WITH AI XAXAXAXAXAXAX
type ColorfulFormatter struct{}

// ANSI codes for backgrounds
const (
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgRed     = "\033[41m"
	bgCyan    = "\033[46m"
	bgMagenta = "\033[45m"
	reset     = "\033[0m"
	gray      = "\033[90m"
)

// Random text colors (foreground)
var randomColors = []string{
	"\033[31m", // red
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[34m", // blue
	"\033[35m", // magenta
	"\033[36m", // cyan
	"\033[37m", // white
}

func (f *ColorfulFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())

	// Choose background based on level
	var bg string
	switch entry.Level {
	case logrus.InfoLevel:
		bg = bgGreen
	case logrus.WarnLevel:
		bg = bgYellow
	case logrus.ErrorLevel:
		bg = bgRed
	case logrus.DebugLevel:
		bg = bgCyan
	case logrus.FatalLevel:
		bg = bgMagenta
	default:
		bg = "\033[47m" // white background fallback
	}

	// Random colorful message
	msgColor := randomColors[rand.Intn(len(randomColors))]
	timeStr := entry.Time.Format("15:04:05")

	// Build log line
	logLine := fmt.Sprintf(
		"%s[%s]%s %s%s%s %s%s%s\n",
		bg, entry.Level.String(), reset,
		gray, timeStr, reset,
		msgColor, entry.Message, reset,
	)

	return []byte(logLine), nil
}

// Hook for logging into a file
type FileHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

// method that must be implemented to be meet AddHook()
func (hook *FileHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&ColorfulFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logFilePath := filepath.Join("logs", "server.log")
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Could not open log file: %v", err)
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&ColorfulFormatter{})

	fileHook := &FileHook{
		Writer:    file,
		LogLevels: logrus.AllLevels,
		Formatter: &logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		},
	}
	logger.AddHook(fileHook)

	return logger
}
