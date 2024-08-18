package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Yuelioi/vidor/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedPath := req.URL.Path

	// Find the index of "/files/"
	index := strings.Index(requestedPath, "/files/")
	if index == -1 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Invalid path"))
		return
	}

	// Get the part after "/files/"
	relativePath := requestedPath[index+len("/files/"):]

	// Convert to normal Windows path
	systemPath := filepath.FromSlash(relativePath)
	println("Requesting file:", systemPath)

	fileData, err := os.ReadFile(systemPath)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", systemPath)))
		return
	}

	res.Write(fileData)
}

func AppLaunch() {
	a := app.Application

	err := wails.Run(&options.App{
		Title:     "Vidor",
		Width:     1050,
		Height:    720,
		MinWidth:  1050,
		MinHeight: 720,

		HideWindowOnClose: true,
		Frameless:         true,
		CSSDragProperty:   "widows",
		CSSDragValue:      "1",
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		LogLevel:           logger.WARNING,
		LogLevelProduction: logger.WARNING,

		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        a.Startup,
		OnShutdown:       a.Shutdown,
		Bind: []interface{}{
			a,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func main() {
	AppLaunch()

	// dl := plugins.NewYTBDownloader()
	// dl.ShowInfo("https://www.youtube.com/watch?v=kaZOXRqFPCw", shared.Config{
	// 	UseProxy: true,
	// 	ProxyURL: "socks://127.0.0.1:10809",
	// })

}
