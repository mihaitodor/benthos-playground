package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/redpanda-data/benthos/v4/public/service"
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"
)

type logger struct {
	label string
	path  string
	slog.Handler
}

func (l *logger) Handle(ctx context.Context, record slog.Record) error {
	if record.Level == slog.LevelError {
		log.Printf("Handling error log record for component %q with path %q", l.label, l.path)
	}
	return l.Handler.Handle(ctx, record)
}

func (l *logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	var label string
	var path string
	var labelFound bool
	var pathFound bool
	for _, a := range attrs {
		if a.Key == "label" {
			label = a.Value.String()
			labelFound = true
		}
		if a.Key == "path" {
			path = a.Value.String()
			pathFound = true
		}
		if labelFound && pathFound {
			break
		}
	}
	if !pathFound {
		// Inherit the path from the parent logger when it's not set
		path = l.path
	}

	return &logger{
		label:   label,
		path:    path,
		Handler: l.Handler.WithAttrs(attrs),
	}
}

func (l *logger) WithGroup(name string) slog.Handler {
	return &logger{
		label:   l.label,
		path:    l.path,
		Handler: l.Handler.WithGroup(name),
	}
}

func main() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	logger := slog.New(&logger{Handler: baseHandler})

	service.RunCLI(context.Background(), service.CLIOptAddTeeLogger(logger))
}
