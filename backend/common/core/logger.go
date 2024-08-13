package core

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"
)

type slogFieldsCtxKey int

type contextHandler struct {
	slog.Handler
}

const slogFieldsCtxID slogFieldsCtxKey = 0

// Handle adds contextual attributes to the Record before calling the underlying.
func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFieldsCtxID).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// WithLogAttr adds an slog attribute to the provided context so that it will be
// included in any Record created with such context.
func WithLogAttr(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFieldsCtxID).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFieldsCtxID, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)

	return context.WithValue(parent, slogFieldsCtxID, v)
}

// newLogger creates a new slog logger with the specified writer and handler options.
// It uses the slog.NewJSONHandler to format log messages as JSON.
func newLogger(writer io.Writer, cfg *slog.HandlerOptions) *slog.Logger {
	return slog.New(&contextHandler{slog.NewJSONHandler(writer, cfg)})
}

// NewNoopLogger returns a logger that discards all log messages.
// It can be used for testing or when no logging is required.
func NewNoopLogger() *slog.Logger {
	return newLogger(io.Discard, nil)
}

// NewStdoutLogger returns a logger that writes log messages to the standard output (os.Stdout).
func NewStdoutLogger() *slog.Logger {
	return newLogger(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				a.Value = slog.AnyValue(slog.Source{
					Function: source.Function,
					File:     filepath.Base(source.File),
					Line:     source.Line,
				})
			}

			return a
		},
	})
}

// NewLoggerModuleWithConfig is an fx.Option that provides a new stdout logger for the application.
// It uses the slog library for structured logging and provides a module for dependency injection.
func NewLoggerModuleWithConfig() fx.Option {
	return fx.Module(
		"Logger Module",
		fx.Provide(NewStdoutLogger),
	)
}
