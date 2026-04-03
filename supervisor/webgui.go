package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed webgui/*
var webFilesRoot embed.FS

// WebFiles is the embedded filesystem with webgui content at root
var WebFiles fs.FS

func init() {
	var err error
	WebFiles, err = fs.Sub(webFilesRoot, "webgui")
	if err != nil {
		panic(err)
	}
}

// SupervisorWebgui the interface to show a WEBGUI to control the supervisor
type SupervisorWebgui struct {
	router     *mux.Router
	supervisor *Supervisor
	pathPrefix string
}

// WebGuiHandler redirects root path to webgui
// The redirect URL includes the path prefix to ensure correct navigation
func (sw *SupervisorWebgui) WebGuiHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := sw.pathPrefix + "/webgui/"
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

// NewSupervisorWebgui create a new SupervisorWebgui object
func NewSupervisorWebgui(supervisor *Supervisor, pathPrefix string) *SupervisorWebgui {
	router := mux.NewRouter()
	return &SupervisorWebgui{router: router, supervisor: supervisor, pathPrefix: pathPrefix}
}

// CreateHandler create a http handler to process the request from WEBGUI
func (sw *SupervisorWebgui) CreateHandler() http.Handler {
	// Redirect root to webgui (e.g., / or /supervisord/ -> /webgui/)
	sw.router.HandleFunc(sw.pathPrefix+"/", sw.WebGuiHandler)

	// Serve static files from embedded filesystem
	// Strip the path prefix and /webgui from the request path
	webguiPath := sw.pathPrefix + "/webgui"
	fileServer := http.FileServer(http.FS(WebFiles))
	sw.router.PathPrefix(webguiPath).Handler(http.StripPrefix(webguiPath, fileServer))

	return sw.router
}
