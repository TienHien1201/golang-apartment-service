package xlogger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zl zerolog.Logger
}

type Config struct {
	Level      string // debug, info, warn, error, fatal, panic
	Format     string // json or console
	Output     string // stdout, stderr, or file path
	TimeFormat string // time format for log messages
}

func New(cfg *Config) (*Logger, error) {
	// Set log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	// Configure output writer
	var output io.Writer
	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("could not open log file: %w", err)
		}
		output = file
	}

	// Configure time format (ensure it's not empty)
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = time.RFC3339Nano
	}
	zerolog.TimeFieldFormat = cfg.TimeFormat

	// If format is "console", use human-readable, otherwise use JSON
	if cfg.Format == "console" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: cfg.TimeFormat,
			NoColor:    false,
		}
	}

	// Create logger instance
	logger := zerolog.New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger()

	return &Logger{zl: logger}, nil
}

// --- Logger methods ---

func (l *Logger) Info(msg string, fields ...Field) {
	event := l.zl.Info()
	for _, field := range fields {
		field.AddTo(event)
	}
	event.Msg(msg)
}

func (l *Logger) Error(msg string, fields ...Field) {
	event := l.zl.Error()
	for _, field := range fields {
		field.AddTo(event)
	}
	event.Msg(msg)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	event := l.zl.Debug()
	for _, field := range fields {
		field.AddTo(event)
	}
	event.Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	event := l.zl.Warn()
	for _, field := range fields {
		field.AddTo(event)
	}
	event.Msg(msg)
}

// Field types for structured logging.
type Field interface {
	AddTo(event *zerolog.Event)
}

type StringField struct {
	Key   string
	Value string
}

func (f StringField) AddTo(event *zerolog.Event) {
	event.Str(f.Key, f.Value)
}

type IntField struct {
	Key   string
	Value int
}

func (f IntField) AddTo(event *zerolog.Event) {
	event.Int(f.Key, f.Value)
}

type Int64Field struct {
	Key   string
	Value int64
}

func (f Int64Field) AddTo(event *zerolog.Event) {
	event.Int64(f.Key, f.Value)
}

type ErrorField struct {
	Key   string
	Value error
}

func (f ErrorField) AddTo(event *zerolog.Event) {
	event.Err(f.Value)
}

type ObjectField struct {
	Key   string
	Value interface{}
}

func (f ObjectField) AddTo(event *zerolog.Event) {
	event.Interface(f.Key, f.Value)
}

// --- Field constructors ---

func String(key, value string) Field {
	return StringField{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return IntField{Key: key, Value: value}
}

func Error(err error) Field {
	return ErrorField{Key: "error", Value: err}
}

func Object(key string, value interface{}) Field {
	return ObjectField{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return IntField{Key: key, Value: int(value / time.Millisecond)}
}

func Int32(key string, value int32) Field {
	return IntField{Key: key, Value: int(value)}
}

func Int64(key string, value int64) Field {
	return Int64Field{Key: key, Value: value}
}

func Strings(key string, value []string) Field {
	return String(key, strings.Join(value, ", "))
}

func Uint(key string, value uint) Field {
	return IntField{Key: key, Value: int(value)}
}

func Uint64(key string, value uint64) Field {
	return Int64Field{Key: key, Value: int64(value)}
}
