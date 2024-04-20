package handlers

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
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
	start := time.Now()
	defer log.Infof("Host: %s\nPath: %s\nTime:%v\n", r.Host, r.RequestURI, time.Since(start))

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		if err := bh.app.Send(r.Method, r.URL.Path, body); err != nil {
			log.Info(err.Error())
			return
		}
	}()

	w.WriteHeader(http.StatusOK)
	return
}
