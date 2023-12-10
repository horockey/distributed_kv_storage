package http_controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (ctrl *httpController) initRouter() {
	ctrl.router = mux.NewRouter()

	ctrl.router.Methods(http.MethodGet).Path("/kv/{key}").HandlerFunc(ctrl.kvGetKey)
	ctrl.router.Methods(http.MethodPost).Path("/kv").HandlerFunc(ctrl.kvPost)

	ctrl.router.HandleFunc("/", ctrl.redirectToDocs)

	ctrl.router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) { w.Write(docsHtml) })
}

func (ctrl *httpController) redirectToDocs(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/docs", http.StatusMovedPermanently)
}
