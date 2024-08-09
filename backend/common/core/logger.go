package core

import (
	"context"
	"io"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

type Logger interface {
	Handler() slog.Handler
	// Debug(msg string, args ...any)
	// DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	// Warn(msg string, args ...any)
	// WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Unwrap() *slog.Logger
}

type LoggerConfig struct {
	RequestHeaderID string
}

type logger struct {
	slog *slog.Logger
	cfg  *LoggerConfig
}

var _ Logger = (*logger)(nil)

func (l logger) Unwrap() *slog.Logger {
	return l.slog
}

func (l logger) Handler() slog.Handler {
	return l.slog.Handler()
}

func (l logger) Info(msg string, args ...any) {
	//nolint:sloglint // Should accept kv-only and attr-only.
	l.slog.Info(msg, args...)
}

func (l logger) InfoContext(ctx context.Context, msg string, args ...any) {
	if l.cfg != nil {
		if requestID, ok := ctx.Value(l.cfg.RequestHeaderID).(string); ok && requestID != "" {
			args = append([]any{
				slog.String("request-id", requestID),
			}, args...)
		}
	}

	//nolint:sloglint // Should accept kv-only and attr-only.
	l.slog.InfoContext(ctx, msg, args...)
}

func (l logger) Error(msg string, args ...any) {
	//nolint:sloglint // Should accept kv-only and attr-only.
	l.slog.Error(msg, args...)
}

func (l logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	if l.cfg != nil {
		if requestID, ok := ctx.Value(l.cfg.RequestHeaderID).(string); ok && requestID != "" {
			args = append([]any{
				slog.String("request-id", requestID),
			}, args...)
		}
	}

	//nolint:sloglint // Should accept kv-only and attr-only.
	l.slog.ErrorContext(ctx, msg, args...)
}

// newLogger creates a new slog logger with the specified writer and handler options.
// It uses the slog.NewJSONHandler to format log messages as JSON.
func newLogger(writer io.Writer, cfg *slog.HandlerOptions) *slog.Logger {
	return slog.New(slog.NewJSONHandler(writer, cfg))
}

// NewNoopLogger returns a logger that discards all log messages.
// It can be used for testing or when no logging is required.
func NewNoopLogger() *logger {
	return &logger{
		slog: newLogger(io.Discard, nil),
	}
}

// NewStdoutLogger returns a logger that writes log messages to the standard output (os.Stdout).
// The logger includes source information and has a debug log level by default.
func NewStdoutLogger(loggerCfg *LoggerConfig) *logger {
	return &logger{
		cfg:  loggerCfg,
		slog: newLogger(os.Stdout, nil),
	}
}

// NewLoggerModuleWithConfig is an fx.Option that provides a new stdout logger for the application.
// It uses the slog library for structured logging and provides a module for dependency injection.
func NewLoggerModuleWithConfig(loggerCfg *LoggerConfig) fx.Option {
	return fx.Module(
		"Logger Module",
		fx.Provide(func() Logger {
			return NewStdoutLogger(loggerCfg)
		}),
	)
}
