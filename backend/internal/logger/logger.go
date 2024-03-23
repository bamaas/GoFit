package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var logLevel map[string]slog.Level = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

// Initialize the logger
func New(l string) (*slog.Logger, error) {

	level, ok := logLevel[strings.ToUpper(l)]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		AddSource: true,
	})), nil

}
