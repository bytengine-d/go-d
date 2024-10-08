package logs

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"log/slog"
	"runtime"
	"strings"
)

const WeShareLogTimeFormat = "2006-01-02 15:04:05,000"
const WeShareLogMsgFormat = "%s [%d] %-7s %s:%d.%s() - %s\n"

type HandlerOptions struct {
	SlogOpts   slog.HandlerOptions
	FilePath   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

type Handler struct {
	slog.Handler
	l *log.Logger
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	timeVal := r.Time.Format(WeShareLogTimeFormat)
	goid := runtime.NumGoroutine()
	upperLevel := strings.ToUpper(r.Level.String())
	file := ""
	funcName := ""
	line := -1
	if r.PC > 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		file = f.File
		if file == "" {
			file = "???"
		}
		funcName = f.Function
		line = f.Line
	}

	h.l.Printf(WeShareLogMsgFormat, timeVal, goid, upperLevel, file, line, funcName, r.Message)
	return nil
}

func NewHandler(
	opts HandlerOptions,
) slog.Handler {
	l := &lumberjack.Logger{
		Filename:   opts.FilePath,
		MaxSize:    opts.MaxSize, // megabytes
		MaxAge:     opts.MaxAge,  //days
		MaxBackups: opts.MaxBackups,
		LocalTime:  opts.LocalTime,
		Compress:   opts.Compress, // disabled by default
	}
	h := &Handler{
		Handler: slog.NewTextHandler(l, &opts.SlogOpts),
		l:       log.New(l, "", 0),
	}
	return h
}
