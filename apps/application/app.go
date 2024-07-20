package application

import (
	"context"
	"github.com/bytengine-d/go-d/space"
	"os"
	"path/filepath"
	"time"
)

const AppNameSpaceKey = "$d_appName"
const AppStartTimeSpaceKey = "$d_appStartTime"

func NewAppFromArgs() context.Context {
	appName := filepath.Base(os.Args[0])
	return NewApp(appName)
}

func NewApp(appName string) context.Context {
	ctx := space.GlobalSpace()
	now := time.Now()
	space.Set(ctx, "appName", appName)
	space.Set(ctx, AppNameSpaceKey, appName)
	space.Set(ctx, AppStartTimeSpaceKey, &now)
	return ctx
}

func GetAppName(ctx context.Context) string {
	if appName, ok := space.GetString(ctx, AppNameSpaceKey); !ok {
		return ""
	} else {
		return appName
	}
}

func GetAppStartTime(ctx context.Context) *time.Time {
	if startupTime, ok := space.Get(ctx, AppStartTimeSpaceKey); !ok {
		return nil
	} else {
		return startupTime.(*time.Time)
	}
}
