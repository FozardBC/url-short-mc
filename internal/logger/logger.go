package logger

import (
	"io"
	"log/slog"
	slogpretty "microservice_t/internal/lib"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelDebug = "debug"
	LevelDev   = "dev"
	LevelProd  = "prod"
)

func New(level string) *slog.Logger {
	var log *slog.Logger

	switch level {
	case LevelDebug:
		log = setupPrettySlog()
	case LevelDev: // TODO: сделать вывод лого куда нибудь во вне
		log = settupDevSlog()
	case LevelProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	}

	return log
}

func settupDevSlog() *slog.Logger {

	fileLogWriter := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileLogWriter)

	return slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // Добавляет файл и строку вызова
	}))
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
