package ioutil

import (
	"context"
)

type WrapperWriterHandler func(writer *WrapperWriter) error

func NewWrapperWithHandler(
	ctx context.Context,
	handlers ...WrapperWriterHandler) (*WrapperWriter, error) {
	wrapper := NewWrapper(ctx)
	var err error
	if handlers != nil && len(handlers) > 0 {
		for _, handler := range handlers {
			err = handler(wrapper)
			if err != nil {
				return nil, err
			}
		}
	}
	return wrapper, nil
}
