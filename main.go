package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title: "GunCode 改枪码管理器",
		// 隐藏系统自带的标题栏/边框，改由前端代码自绘标题栏与窗口控制按钮
		Frameless: true,
		// 这里的尺寸是逻辑像素，Windows 会按显示器缩放比例换算成物理像素
		// （150% 缩放时 1180 逻辑像素 = 1770 物理像素）。
		Width:     1180,
		Height:    760,
		MinWidth:  1024,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 21, G: 23, B: 28, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
