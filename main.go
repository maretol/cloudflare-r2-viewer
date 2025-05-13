package main

import (
	"cloudflare-r2-viewer/backend"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()
	viewerHandler := backend.NewViewerHandler()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "cloudflare-r2-viewer",
		Width:             1200,
		Height:            800,
		MinWidth:          1024,
		MinHeight:         768,
		MaxWidth:          1920,
		MaxHeight:         1080,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         logger.DEBUG,
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		OnShutdown:       app.shutdown,
		WindowStartState: options.Normal,
		Bind: []any{
			app,
			viewerHandler,
		},
		// Windows platform specific options
		Windows: windowsConfig(),
		// Mac platform specific options
		Mac: macConfig(),
	})

	if err != nil {
		log.Fatal(err)
	}
}

func windowsConfig() *windows.Options {
	return &windows.Options{
		WebviewIsTransparent: false,
		WindowIsTranslucent:  false,
		DisableWindowIcon:    false,
		// DisableFramelessWindowDecorations: false,
		WebviewUserDataPath: "",
		ZoomFactor:          1.0,
	}
}

func macConfig() *mac.Options {
	return &mac.Options{
		TitleBar:             nil,
		Appearance:           mac.NSAppearanceNameDarkAqua,
		WebviewIsTransparent: true,
		WindowIsTranslucent:  true,
		About: &mac.AboutInfo{
			Title:   "cloudflare-r2-viewer",
			Message: "",
			Icon:    icon,
		},
	}
}
