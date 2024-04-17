package handlers

import (
	"io"
	"net/http"
)

type App interface {
	Send(string, string, []byte) error
}

type BalancingHandler struct {
	app App
}

func New(app App) *BalancingHandler {
	return &BalancingHandler{
		app: app,
	}
}

func (bh *BalancingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bh.app.Send(r.Method, r.URL.String(), body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
