package finder

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var failurePercentage = 0.5

type DB interface {
	Set(context.Context, models.Host) error
}

type Finder struct {
	foundServices   *sync.Map
	db              DB
	baseServiceName string
	portToScan      int
}

func New(m *sync.Map, baseName string, port int, db DB) *Finder {
	return &Finder{
		foundServices:   m,
		baseServiceName: baseName,
		portToScan:      port,
		db:              db,
	}
}

func (f *Finder) Scan() {
	fails := 0
	limitMet := false
	n := 0

	for !limitMet {
		n++

		key := fmt.Sprintf("%s-%d:%d", f.baseServiceName, n, f.portToScan)

		resp, err := http.Get(fmt.Sprintf("http://%s-%d:%d/version", f.baseServiceName, n, f.portToScan))

		if err != nil || resp.StatusCode != 200 {
			fails++
			f.foundServices.Delete(key)
			if float64(n/fails) > failurePercentage {
				limitMet = true
			}
			continue
		}

		log.WithField("key", key).Info("found service")

		if _, ok := f.foundServices.Load(key); !ok {
			f.foundServices.Store(key, struct{}{})
			log.WithFields(log.Fields{
				"service": f.baseServiceName,
				"n":       n,
				"port":    f.portToScan,
			}).Info("storing host")
			if err := f.db.Set(context.Background(), models.Host{
				Address: fmt.Sprintf("%s-%d", f.baseServiceName, n),
				Port:    f.portToScan,
			}); err != nil {
				log.WithField("err", err).Error("could not set host info to db")
				limitMet = true
			}
		}
	}
}
