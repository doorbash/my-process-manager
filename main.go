package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	var runtimeContext context.Context
	var systray *application.SystemTray

	mainApp := application.NewWithOptions(&options.App{
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
		OnStartup: func(ctx context.Context) {
			runtimeContext = ctx
			app.startup(runtimeContext)
		},
		OnDomReady: app.domReady,
		OnShutdown: app.shutdown,
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

	var showWindow = func() {
		// Show the window
		// In a future version of this API, it will be possible to
		// create windows programmatically and be able to show/hide
		// them from the systray with something like:
		//
		// myWindow := mainApp.NewWindow(...)
		// mainApp.NewSystemTray(&options.SystemTray{
		//   OnLeftClick: func() {
		//      myWindow.SetVisibility(!myWindow.IsVisible())
		//   }
		// })
		runtime.Show(runtimeContext)
	}

	systray = mainApp.NewSystemTray(&options.SystemTray{
		// This is the icon used when the system in using light mode
		LightModeIcon: &options.SystemTrayIcon{
			Data: Icon,
		},
		// This is the icon used when the system in using dark mode
		DarkModeIcon: &options.SystemTrayIcon{
			Data: Icon,
		},
		Tooltip:     APP_TITLE,
		OnLeftClick: showWindow,
		OnMenuClose: func() {
			// Add the left click call after 500ms
			// We do this because the left click fires right
			// after the menu closes, and we don't want to show
			// the window on menu close.
			go func() {
				time.Sleep(500 * time.Millisecond)
				systray.OnLeftClick(showWindow)
			}()
		},
		OnMenuOpen: func() {
			// Remove the left click callback
			systray.OnLeftClick(func() {})
		},
	})

	systray.SetMenu(menu.NewMenuFromItems(
		menu.Label(fmt.Sprintf("show %s", APP_TITLE)).OnClick(func(c *menu.CallbackData) {
			showWindow()
		}),
		menu.Label("quit").OnClick(func(_ *menu.CallbackData) {
			mainApp.Quit()
		}),
	))

	if err := mainApp.Run(); err != nil {
		log.Fatal(err)
	}
}
