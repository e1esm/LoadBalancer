package db

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type HostsDB struct {
	cli *redis.Client
}

func New(addr string, password string, db int) *HostsDB {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})

	if cli.Ping(context.Background()).Err() != nil {
		log.Fatal("could not connect to redis")
	}
	return &HostsDB{
		cli: cli,
	}
}

func (h *HostsDB) GetHosts(ctx context.Context) ([]models.Host, error) {
	cmd := h.cli.Do(ctx, "KEYS", "*")

	res, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("error while fetching values: %w", err)
	}

	if hosts := res.([]models.Host); hosts != nil {
		return hosts, nil
	}

	return nil, fmt.Errorf("parsing result error")
}

func (h *HostsDB) Set(ctx context.Context, host models.Host) error {
	return h.cli.Set(ctx, host.Address, host, 0).Err()
}

func (h *HostsDB) Get(ctx context.Context, addr string) (models.Host, error) {
	var host models.Host
	err := h.cli.Get(ctx, addr).Scan(&host)

	return host, fmt.Errorf("scanning error: %w", err)
}

func (h *HostsDB) Close() {
	err := h.cli.Close()
	if err != nil {
		log.WithError(err).Error("database connection closing error")
	}
}
