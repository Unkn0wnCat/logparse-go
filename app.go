package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Cowabunga",
		Filters: []runtime.FileFilter{
			runtime.FileFilter{
				DisplayName: "Log Files",
				Pattern:     "*.log",
			},
		},
		ShowHiddenFiles:            false,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Hello %s, It's show time!", path)
}
