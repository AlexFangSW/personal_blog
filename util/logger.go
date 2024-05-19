package util

import (
	"log/slog"
)

var levelMap = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
}

func InitLogger(level string) {
	realLevel, ok := levelMap[level]
	if !ok {
		realLevel = slog.LevelInfo
		slog.Error("Provided log level doesn't exist, default to: INFO", "provided", level)
	}

	slog.SetLogLoggerLevel(realLevel)
}
