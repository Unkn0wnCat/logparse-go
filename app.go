package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"logparse-go/database"
	"logparse-go/importer"
	"logparse-go/parser"
	"logparse-go/resultcollector"
	"strings"
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

func (a *App) RunImport(logFile string, databaseFile string) ([]resultcollector.Result, error) {
	collector, err := importer.Import(logFile, databaseFile)
	if err != nil {
		return nil, err
	}

	return collector.GetAll(), nil
}

func (a *App) QueryView(databaseFile string, ipFilter string, timestampFilterFrom string, timestampFilterTo string) ([]parser.LogLine, error) {
	db, err := database.OpenDatabase(databaseFile)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	return db.Query(ipFilter, strings.ReplaceAll(timestampFilterFrom, "T", " "), strings.ReplaceAll(timestampFilterTo, "T", " "))
}

func (a *App) QueryIPCounts(databaseFile string, ipFilter string, timestampFilterFrom string, timestampFilterTo string) ([]database.IPCount, error) {
	db, err := database.OpenDatabase(databaseFile)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	return db.QueryIPCounts(ipFilter, strings.ReplaceAll(timestampFilterFrom, "T", " "), strings.ReplaceAll(timestampFilterTo, "T", " "))
}

func (a *App) QueryStatusCounts(databaseFile string, ipFilter string, timestampFilterFrom string, timestampFilterTo string, codes string) ([]database.StatusCount, error) {
	db, err := database.OpenDatabase(databaseFile)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	return db.QueryStatusCounts(ipFilter, strings.ReplaceAll(timestampFilterFrom, "T", " "), strings.ReplaceAll(timestampFilterTo, "T", " "), codes)
}

func (a *App) PickLogFile() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Log File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Log Files",
				Pattern:     "*.log",
			},
		},
		ShowHiddenFiles:            true,
		CanCreateDirectories:       false,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return err.Error()
	}
	return path
}

func (a *App) PickDBFile() string {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "logparse.db",
		Title:           "Select Database File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Database Files",
				Pattern:     "*.db",
			},
		},
		ShowHiddenFiles:            true,
		CanCreateDirectories:       false,
		TreatPackagesAsDirectories: false,
	})
	if err != nil {
		return err.Error()
	}
	return path
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
