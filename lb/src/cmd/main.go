package main

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/app"
	"github.com/e1esm/LoadBalancer/lb/src/balancer"
	"github.com/e1esm/LoadBalancer/lb/src/cmd/config"
	"github.com/e1esm/LoadBalancer/lb/src/db"
	"github.com/e1esm/LoadBalancer/lb/src/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {

	cfg := config.New()

	database := db.New(cfg.Redis.Address, cfg.Redis.Password, cfg.Redis.DB)

	intervalDuration, err := time.ParseDuration(cfg.ResetInterval)
	if err != nil {
		intervalDuration = 1 * time.Minute
		log.WithError(err).Error("duration was not parsed")
	}

	blnc := balancer.New(database, cfg.MaxCapacity, intervalDuration)

	appl := app.New(blnc)

	hndl := handlers.New(appl)

	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), hndl))
	}()

	go func() {
		blnc.Reset(ctx)
	}()

	select {
	case <-exit:
		cancel()
		database.Close()
	}
}
