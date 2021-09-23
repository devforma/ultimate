package web

import (
	"log"
	"net/http"

	"github.com/devforma/ultimate/internal/util"
)

type handler struct {
	notFoundContent string
	routes          map[string]http.HandlerFunc
}

func NewHandler(notFoundContent string) *handler {
	return &handler{
		notFoundContent: notFoundContent,
		routes:          make(map[string]http.HandlerFunc),
	}
}

func (h *handler) AddRoute(method string, path string, fn http.HandlerFunc) {
	h.routes[method+":"+path] = fn
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	fn, ok := h.routes[r.Method+":"+r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(util.StringToBytes(h.notFoundContent))
		return
	}

	fn(w, r)
}
