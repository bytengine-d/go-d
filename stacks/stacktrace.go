package stacks

import (
	"github.com/bytengine-d/go-d/buffer"
	"github.com/bytengine-d/go-d/pool"
	"runtime"
)

// region fields

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)

var _stacktracePool = pool.New(func() *stacktrace {
	return &stacktrace{
		storage: make([]uintptr, 64),
	}
})

type stacktraceDepth int

const (
	stacktraceFirst stacktraceDepth = iota

	stacktraceFull
)

// endregion

// region struct
type stacktrace struct {
	pcs    []uintptr
	frames *runtime.Frames

	storage []uintptr
}

func (st *stacktrace) Free() {
	st.frames = nil
	st.pcs = nil
	_stacktracePool.Put(st)
}

func (st *stacktrace) Count() int {
	return len(st.pcs)
}

func (st *stacktrace) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

type stackFormatter struct {
	b        *buffer.Buffer
	nonEmpty bool
}

func newStackFormatter(b *buffer.Buffer) stackFormatter {
	return stackFormatter{b: b}
}

func (sf *stackFormatter) FormatStack(stack *stacktrace) {
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

func (sf *stackFormatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.AppendByte('\n')
	}
	sf.nonEmpty = true

	sf.b.AppendString("\tat ")
	sf.b.AppendString(frame.File)
	sf.b.AppendByte(':')
	sf.b.AppendInt(int64(frame.Line))
	sf.b.AppendByte('\t')
	sf.b.AppendString(frame.Function)
}

// endregion

// region functions

// region private functions
func captureStacktrace(skip int, depth stacktraceDepth) *stacktrace {
	stack := _stacktracePool.Get()

	switch depth {
	case stacktraceFirst:
		stack.pcs = stack.storage[:1]
	case stacktraceFull:
		stack.pcs = stack.storage
	}

	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	if depth == stacktraceFull {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

func takeStacktrace(frames *runtime.Frames) string {
	buf := Get()
	defer buf.Free()

	stackFmt := newStackFormatter(buf)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		stackFmt.FormatFrame(frame)
	}
	return buf.String()
}

// endregion

// region public functions
func TakeStacktrace(skip int) string {
	stack := captureStacktrace(skip+1, stacktraceFull)
	defer stack.Free()

	return takeStacktrace(stack.frames)
}

func Frames(skip int) *runtime.Frames {
	stack := captureStacktrace(skip+1, stacktraceFull)
	defer stack.Free()
	return stack.frames
}

func StringFrames(frames *runtime.Frames) string {
	return takeStacktrace(frames)
}

// endregion

// endregion
