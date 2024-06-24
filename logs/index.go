package logs

import (
	"context"
	"github.com/bytengine-d/go-d/ioutil"
	"io"
	"log/slog"
)

var root *slog.Logger
var rootWrapperWriter *ioutil.WrapperWriter

func SetupStd(ctx context.Context, slogOpts slog.HandlerOptions) {
	rootWrapperWriter = ioutil.NewWrapperWithStdio(ctx)
	root = slog.New(NewWeShareHandler(rootWrapperWriter, WeShareHandlerOptions{
		SlogOpts: slogOpts,
	}))
	slog.SetDefault(root)
}

func SetupWithFile(ctx context.Context, slogOpts WeShareHandlerOptions) error {
	var err error
	rootWrapperWriter, err = ioutil.NewWrapperWithHandler(ctx, ioutil.DefaultRollDay(slogOpts.FilePath))
	if err != nil {
		return err
	}
	root = slog.New(NewWeShareHandler(rootWrapperWriter, slogOpts))
	slog.SetDefault(root)
	return nil
}

func SetRootWriter(writer io.Writer) {
	rootWrapperWriter.ChangeDelegate(writer)
}
