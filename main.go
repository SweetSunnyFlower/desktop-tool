package main

import (
	"embed"

	appConfig "tools/config"
	"tools/pkg/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed .env
var envContent []byte

func init() {
	// 加载 config 目录下的配置信息
	appConfig.Initialize()
}

func main() {
	// 加载配置文件
	config.InitConfig(envContent)

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  config.GetString("app.name"),
		Width:  config.GetInt("app.width"),
		Height: config.GetInt("app.height"),
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
