package main

import "github.com/julienschmidt/httprouter"

func router(r *httprouter.Router, hdl *handler) {
	r.GET("/", hdl.home)
}