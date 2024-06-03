package util

import (
	"context"
	"fmt"
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

func LogQuery(ctx context.Context, prefix, stmt string) {
	if fire := slog.Default().Enabled(ctx, slog.LevelDebug); fire {
		fmt.Println(prefix, stmt)
	}
}
