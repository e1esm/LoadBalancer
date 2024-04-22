package balancer

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type DB interface {
	Set(context.Context, models.Host) error
	GetOrSet(context.Context, string) (models.Host, error)
}

type LoadBalancer struct {
	maxCapacity       int
	resetInterval     time.Duration
	db                DB
	availableServices *sync.Map
}

func New(db DB, max int, interval time.Duration, services *sync.Map) *LoadBalancer {
	return &LoadBalancer{
		db:                db,
		maxCapacity:       max,
		resetInterval:     interval,
		availableServices: services,
	}
}

func (lb *LoadBalancer) DropHost(ctx context.Context, address string) error {
	h, err := lb.db.GetOrSet(ctx, address)
	if err != nil {
		return fmt.Errorf("dropping host error: %w", err)
	}
	h.Stats.OngoingReqs--

	return lb.db.Set(context.Background(), h)
}

func (lb *LoadBalancer) AcquireHost(ctx context.Context) (*Host, error) {
	var host Host
	minOngoingReqs := lb.maxCapacity
	hosts := lb.availableHosts(ctx)
	var pickedHostIndex int

	for i, h := range hosts {
		if h.Stats.OngoingReqs < minOngoingReqs {
			minOngoingReqs = h.Stats.OngoingReqs

			host = Host{
				ip:   h.Address,
				port: h.Port,
				n:    i + 1,
			}
			pickedHostIndex = i
		}
	}

	if len(hosts) > 0 {
		hosts[pickedHostIndex].Stats.OngoingReqs++
		hosts[pickedHostIndex].Stats.OverallRequests++

		log.WithField("hosts", hosts).WithField("host", host).Info("updating stats")

		if err := lb.db.Set(ctx, hosts[pickedHostIndex]); err != nil {
			return nil, fmt.Errorf("stats udpdate error: %w", err)
		}
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

	hosts := lb.availableHosts(ctx)

	for _, host := range hosts {
		host.Stats.OngoingReqs = 0

		if err := lb.db.Set(ctx, host); err != nil {
			return fmt.Errorf("resetting current capacity error: %w", err)
		}
	}

	return nil
}

func (lb *LoadBalancer) setHostInfo(ctx context.Context, addr ...string) error {
	if err := lb.db.Set(ctx, models.Host{
		Address: addr[0],
		Port:    strToInt(addr[1]),
	}); err != nil {
		return fmt.Errorf("error saving host")
	}

	return nil
}
