package logs

import (
	"context"
	"io"
	"log"
	"log/slog"
	"runtime"
	"strings"
)

const WeShareLogTimeFormat = "2006-01-02 15:04:05,000"
const WeShareLogMsgFormat = "%s [%d] %-7s %s:%d.%s() - %s\n"

type WeShareHandlerOptions struct {
	SlogOpts slog.HandlerOptions
	FilePath string
}

type WeShareHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *WeShareHandler) Handle(ctx context.Context, r slog.Record) error {
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

func NewWeShareHandler(
	out io.Writer,
	opts WeShareHandlerOptions,
) slog.Handler {
	h := &WeShareHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}
