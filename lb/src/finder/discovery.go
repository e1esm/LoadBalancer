package finder

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type Finder struct {
	foundServices   *sync.Map
	baseServiceName string
	portToScan      int
}

func New(m *sync.Map, baseName string, port int) *Finder {
	return &Finder{
		foundServices:   m,
		baseServiceName: baseName,
		portToScan:      port,
	}
}

func (f *Finder) Scan(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.scan(ctx)
		}
	}
}

func (f *Finder) scan(ctx context.Context) {
	limitMet := false
	n := 1
	client := http.DefaultClient

	for !limitMet {
		key := fmt.Sprintf("%s-%d:%d", f.baseServiceName, n, f.portToScan)

		req, err := http.NewRequestWithContext(ctx,
			"GET",
			fmt.Sprintf("http://%s-%d:%d/version", f.baseServiceName, n, f.portToScan),
			nil)

		if err != nil {
			log.WithError(err).Error("error while preparing healthcheck request")
		}

		_, err = client.Do(req)
		if err != nil {
			limitMet = true
			continue
		}

		f.foundServices.Store(key, struct{}{})
		n++
	}
}
