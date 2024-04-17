package app

import (
	"context"
	"github.com/e1esm/LoadBalancer/lb/src/balancer"
)

type Balancer interface {
	ReceiverHost(context.Context) (*balancer.Host, error)
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
	return nil
}
