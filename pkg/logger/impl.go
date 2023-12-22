package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Impl struct {
	log *slog.Logger

	service string
	sentry  *sentry.Client
}

type Opts struct {
	Env string

	Sentry *sentry.Client
	Level  slog.Level
}

var _ Logger = (*Impl)(nil)

func New(opts Opts) *Impl {
	level := opts.Level

	var zeroLogWriter io.Writer
	if opts.Env == "production" {
		zeroLogWriter = os.Stderr
	} else {
		zeroLogWriter = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	slogzerolog.SourceKey = "source"
	slogzerolog.ErrorKeys = []string{"error", "err"}
	zerolog.ErrorStackFieldName = "stack"

	zeroLogLogger := zerolog.New(zeroLogWriter)

	log := slog.New(
		slogmulti.Fanout(
			slogzerolog.Option{
				Level:     level,
				Logger:    &zeroLogLogger,
				AddSource: false,
			}.NewZerologHandler(),
			slogsentry.Option{Level: slog.LevelError, AddSource: true}.NewSentryHandler(),
		),
	)

	return &Impl{
		log:    log,
		sentry: opts.Sentry,
	}
}

func (c *Impl) handle(level slog.Level, input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, input, pcs[0])
	for _, f := range fields {
		r.Add(f)
	}
	_ = c.log.Handler().Handle(context.Background(), r)
}

func (c *Impl) Info(input string, fields ...any) {
	c.log.Info(input, fields...)
}

func (c *Impl) Warn(input string, fields ...any) {
	c.log.Warn(input, fields...)
}

func (c *Impl) Error(input string, fields ...any) {
	c.log.Error(input, fields...)
}

func (c *Impl) Debug(input string, fields ...any) {
	c.log.Debug(input, fields...)
}

func (c *Impl) WithComponent(name string) Logger {
	return &Impl{
		log:     c.log.With(slog.String("component", name)),
		sentry:  c.sentry,
		service: c.service,
	}
}

func (c *Impl) GetSlog() *slog.Logger {
	return c.log
}
