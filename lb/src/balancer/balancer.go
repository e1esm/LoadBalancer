package balancer

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	log "github.com/sirupsen/logrus"
	"time"
)

type DB interface {
	GetHosts(context.Context) ([]models.Host, error)
	Set(context.Context, models.Host) error
	Get(context.Context, string) (models.Host, error)
}

type LoadBalancer struct {
	maxCapacity   int
	resetInterval time.Duration
	db            DB
}

func New(db DB, max int, interval time.Duration) *LoadBalancer {
	return &LoadBalancer{
		db:            db,
		maxCapacity:   max,
		resetInterval: interval,
	}
}

func (lb *LoadBalancer) DropHost(ctx context.Context, address string) error {
	h, err := lb.db.Get(ctx, address)
	if err != nil {
		return fmt.Errorf("dropping host error: %w", err)
	}
	h.Stats.OngoingReqs--

	return lb.db.Set(context.Background(), h)
}

func (lb *LoadBalancer) AcquireHost(ctx context.Context) (*Host, error) {
	var host Host

	minOngoingReqs := lb.maxCapacity

	hosts, err := lb.db.GetHosts(ctx)

	if err != nil {
		return nil, fmt.Errorf("no host was found: %w", err)
	}

	var pickedHostIndex int

	for i, h := range hosts {
		if h.Stats.OngoingReqs < minOngoingReqs {
			minOngoingReqs = h.Stats.OngoingReqs

			host = Host{
				ip:   h.Address,
				port: h.Port,
			}
			pickedHostIndex = i
		}
	}

	hosts[pickedHostIndex].Stats.OngoingReqs++
	hosts[pickedHostIndex].Stats.OverallRequests++

	if err := lb.db.Set(ctx, hosts[pickedHostIndex]); err != nil {
		return nil, fmt.Errorf("stats udpdate error: %w", err)
	}

	return &host, nil
}

func (lb *LoadBalancer) Reset(ctx context.Context) {
	t := time.NewTicker(lb.resetInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := lb.reset(ctx); err != nil {
				log.WithError(err).Error("reset error")
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (lb *LoadBalancer) reset(ctx context.Context) error {

	hosts, err := lb.db.GetHosts(ctx)
	if err != nil {
		return fmt.Errorf("reset error: %w", err)
	}

	for _, host := range hosts {
		host.Stats.OngoingReqs = 0

		if err := lb.db.Set(ctx, host); err != nil {
			return fmt.Errorf("resetting current capacity error: %w", err)
		}
	}

	return nil
}
