package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

type Server struct {
	logger *log.Logger
}

func New() *Server {
	return &Server{
		logger: log.New(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.String() == "/version" {
		w.WriteHeader(http.StatusOK)
		return
	}

	toSleep := time.Duration(rand.Intn(2))
	time.Sleep(toSleep * time.Second)

	start := time.Now()
	defer s.logger.Info(time.Since(start))

	log.Infof("Requested URL: %s", r.URL.Path)
	log.Infof("Requesting host: %s", r.Host)

	w.WriteHeader(http.StatusOK)
}
