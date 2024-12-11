package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"logparse-go/parser"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	file, err := os.OpenFile("./netnea sample-access.logchangedIPs.log", os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(parser.ParseLine(scanner.Text()))
		break
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Log Parser",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
