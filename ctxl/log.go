package ctxl

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"maps"
)

const KeyContextLoggerMeta = "KeyContextLoggerMeta"
const KeyLoggerName = "logger_name"

var emptyMeta = map[string]interface{}{}

type Logger func(ctx context.Context) *zerolog.Logger

func New(name string) func(ctx context.Context) *zerolog.Logger {
	return func(ctx context.Context) *zerolog.Logger {
		meta := GetFields(ctx)
		logger := log.With().Stack().Str(KeyLoggerName, name).Fields(meta).Logger()
		return &logger
	}
}

func SetFields(ctx context.Context, values map[string]interface{}) context.Context {
	orig := GetFields(ctx)
	maps.Copy(orig, values)
	return context.WithValue(ctx, KeyContextLoggerMeta, orig)
}

func SetField(ctx context.Context, key string, value interface{}) context.Context {
	valuesAny := ctx.Value(KeyContextLoggerMeta)

	values, ok := valuesAny.(map[string]interface{})
	if valuesAny != nil && !ok {
		log.
			Err(fmt.Errorf("context already containes %s key but it is not expected type map[string]interface{}")).
			Any(KeyContextLoggerMeta, values).
			Msg("SetMetaField")
		return ctx
	}

	if values == nil {
		values = map[string]interface{}{
			key: value,
		}
	} else {
		values[key] = value
	}

	return context.WithValue(ctx, KeyContextLoggerMeta, values)
}

func GetFields(ctx context.Context) map[string]interface{} {
	metaAny := ctx.Value(KeyContextLoggerMeta)
	if metaAny == nil {
		return emptyMeta
	}

	meta, ok := metaAny.(map[string]interface{})
	if !ok {
		log.
			Err(fmt.Errorf("context already contains %s key but it is not expected type map[string]interface{}", KeyContextLoggerMeta)).
			Any(KeyContextLoggerMeta, meta).
			Msg("GetMeta")
		return map[string]interface{}{}
	}

	return meta
}
