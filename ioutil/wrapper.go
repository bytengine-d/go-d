package ioutil

import "io"

type WrapperWriter struct {
	delegate io.Writer
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

func Wrapper(writer io.Writer) *WrapperWriter {
	return &WrapperWriter{delegate: writer}
}
