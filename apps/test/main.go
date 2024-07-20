package main

import (
	"fmt"
	"github.com/bytengine-d/go-d/apps/application"
	"github.com/bytengine-d/go-d/logs"
	"log"
	"log/slog"
)

func main() {
	ctx := application.NewAppFromArgs()
	err := logs.SetupWithFile(ctx, logs.WeShareHandlerOptions{
		FilePath: fmt.Sprintf("tmp/log/%s.log", application.GetAppName(ctx)),
	})
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("app launched")
}
