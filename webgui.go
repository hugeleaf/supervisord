package main

import (
	"net/http"
	"embed"
	"github.com/gorilla/mux"
)

//go:embed webgui/*
var WebFlies embed.FS

// SupervisorWebgui the interface to show a WEBGUI to control the supervisor
type SupervisorWebgui struct {
	router     *mux.Router
	supervisor *Supervisor
}

// NewSupervisorWebgui create a new SupervisorWebgui object
func NewSupervisorWebgui(supervisor *Supervisor) *SupervisorWebgui {
	router := mux.NewRouter()
	return &SupervisorWebgui{router: router, supervisor: supervisor}
}

func WebGuiHandler(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/webgui", http.StatusMovedPermanently)
}

// CreateHandler create a http handler to process the request from WEBGUI
func (sw *SupervisorWebgui) CreateHandler() http.Handler {
	// sw.router.PathPrefix("/").Handler(http.FileServer(HTTP))
    sw.router.HandleFunc("/", WebGuiHandler)
	sw.router.PathPrefix("/webgui").Handler(http.FileServer(http.FS(WebFlies)))
	return sw.router
}
