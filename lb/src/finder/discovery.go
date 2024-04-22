package finder

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
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

func (f *Finder) Scan() {

	limitMet := false
	n := 1
	for !limitMet {
		key := fmt.Sprintf("%s-%d:%d", f.baseServiceName, n, f.portToScan)

		resp, err := http.Get(fmt.Sprintf("http://%s-%d:%d/version", f.baseServiceName, n, f.portToScan))
		if err != nil {
			log.WithError(err).Info("Http request has been made to ser")
		}

		if resp.StatusCode != 200 {
			limitMet = true
			f.foundServices.Delete(key)
			continue
		}

		if _, ok := f.foundServices.Load(key); !ok {
			f.foundServices.Store(key, struct{}{})
		}

		n++
	}
}
