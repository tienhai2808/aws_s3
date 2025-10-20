package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
}

func newHandler() *handler {
	return &handler{}
}

func (h *handler) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Welcome to the Home Page.")
}
