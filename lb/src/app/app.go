package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/balancer"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	httpTimeout = 5 * time.Second
)

type Balancer interface {
	AcquireHost(context.Context) (*balancer.Host, error)
	DropHost(ctx context.Context, address string) error
}

type App struct {
	balancer Balancer
}

func New(balancer Balancer) *App {
	return &App{
		balancer: balancer,
	}
}

func (a *App) Send(method string, url string, body []byte) error {
	host, err := a.balancer.AcquireHost(context.Background())
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	defer a.balancer.DropHost(context.Background(), host.Address())

	req := request{
		method: method,
		url:    "http://" + host.Address() + url,
		body:   body,
	}

	if err = a.send(req); err != nil {
		log.WithError(err).WithField("req", req).Error("error while sending http request to target")
		return err
	}

	return nil
}

func (a *App) send(req request) error {
	cli := http.Client{Timeout: httpTimeout}
	var (
		resp *http.Response
		err  error
	)

	switch req.method {
	case http.MethodGet:
		resp, err = cli.Get(req.url)
	case http.MethodPost:
		resp, err = cli.Post(req.url, "application/json", bytes.NewBuffer(req.body))
	default:
		return fmt.Errorf("unsupported HTTP request")
	}

	if err != nil {
		return fmt.Errorf("http request error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failure")
	}

	return nil
}
