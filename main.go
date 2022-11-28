package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed frontend/build
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

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
		Width:             480,
		Height:            800,
		MinWidth:          480,
		MinHeight:         800,
		MaxWidth:          1280,
		MaxHeight:         740,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: true,
		Assets:            assets,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
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
	})

	if err != nil {
		log.Fatal(err)
	}

}
