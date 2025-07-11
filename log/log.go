package logconfig

import (
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	EnvLogMode = "MYSTERIUM_LOG_MODE"
	ModeJSON   = "json"
)

type Config struct {
	LogFilePath string
	LogFileName string
}

func BootstrapDefaultLogger(writers ...io.Writer) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	if isJSONMode() {
		writers = append(writers, os.Stderr)
	} else {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := log.Output(zerolog.MultiLevelWriter(writers...)).
		Level(zerolog.DebugLevel).
		With().
		Caller().
		Timestamp().
		Logger()

	setGlobalLogger(&logger)

	return &logger
}

// SetLevel given a logger sets it's level if possible.
// If a given level string is not parseable, default log level is used.
func SetLevel(l *zerolog.Logger, level string) *zerolog.Logger {
	setLevel := zerolog.DebugLevel
	if lvl, err := zerolog.ParseLevel(level); err == nil {
		setLevel = lvl
	}

	logger := l.Level(setLevel)
	setGlobalLogger(&logger)

	return &logger
}

func setGlobalLogger(logger *zerolog.Logger) {
	log.Logger = *logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

func isJSONMode() bool {
	v, ok := os.LookupEnv(EnvLogMode)
	if !ok {
		return false
	}
	return v == ModeJSON
}
