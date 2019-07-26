package api

import (
	"github.com/gorilla/mux"
)

func (a *App) endpoints() {
	a.Router.HandleFunc("/auth", a.login)
	a.Router.HandleFunc("/warn", a.restricted(notImplemented))
	a.Router.HandleFunc("/warn/{id}", a.restricted(notImplemented))
	a.Router.HandleFunc("/mails", a.restricted(a.mailGeneral))
	a.Router.HandleFunc("/mails/add", a.restricted(a.mailRegister))
	a.Router.HandleFunc("/mails/{id}/delete", a.restricted(a.mailDelete))
	a.Router.HandleFunc("/mails/{id}/update", a.restricted(a.mailEdit))
	a.Router.HandleFunc("/upload", a.restricted(a.uploadCSV))
}

//Initialize the router
func Initialize() *App {
	thisRouter := mux.NewRouter()
	thisApp := &App{Router: thisRouter}
	thisApp.initKey()
	thisApp.endpoints()
	return thisApp
}
