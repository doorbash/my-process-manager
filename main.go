package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var runtimeContext context.Context

const (
	GITHUB_URL = "https://github.com/doorbash/my-process-manager"
	DB_NAME    = "db.sqlite"
	APP_TITLE  = "My Process Manager"
)

func main() {
	app := NewApp(
		GITHUB_URL,
		DB_NAME,
	)

	err := wails.Run(&options.App{
		Title:             APP_TITLE,
		Width:             500,
		Height:            800,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		LogLevel:          logger.DEBUG,
		OnStartup: func(ctx context.Context) {
			runtimeContext = ctx
			app.startup(runtimeContext)
		},
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: APP_TITLE,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   APP_TITLE,
				Message: "",
				Icon:    icon,
			},
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
