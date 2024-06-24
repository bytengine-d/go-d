package ioutil

import (
	"context"
	"github.com/bytengine-d/go-d/space"
	"io"
	"os"
)

type WrapperWriter struct {
	delegate io.Writer
	ctx      context.Context
}

func (w *WrapperWriter) Write(p []byte) (n int, err error) {
	return w.delegate.Write(p)
}

func (w *WrapperWriter) Delegate() io.Writer {
	return w.delegate
}

func (w *WrapperWriter) ChangeDelegate(writer io.Writer) {
	w.delegate = writer
}

func (w *WrapperWriter) Set(key string, val any) *WrapperWriter {
	space.Set(w.ctx, key, val)
	return w
}

func (w *WrapperWriter) Has(key string) bool {
	return space.Has(w.ctx, key)
}

func (w *WrapperWriter) Get(key string) (any, bool) {
	return space.Get(w.ctx, key)
}

func (w *WrapperWriter) Delete(key string) *WrapperWriter {
	space.Remove(w.ctx, key)
	return w
}

func NewWrapper(ctx context.Context) *WrapperWriter {
	return &WrapperWriter{ctx: ctx}
}

func NewWrapperWithStdio(ctx context.Context) *WrapperWriter {
	return &WrapperWriter{
		ctx:      ctx,
		delegate: os.Stdout,
	}
}

func NewWrapperWithDelegate(ctx context.Context, writer io.Writer) *WrapperWriter {
	wrapper := &WrapperWriter{
		delegate: writer,
		ctx:      space.NewSpace(ctx, "WrapperWriterSpace"),
	}

	return wrapper
}
