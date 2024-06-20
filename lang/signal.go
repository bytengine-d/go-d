package lang

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

func AddSignalHandler(handler func() int) {
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		result := -1
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT:
				result = handler()
			default:
				slog.Info(fmt.Sprintf("其他信号: %s", s))
			}

			if result > -1 {
				break
			}
		}
	}()
}
