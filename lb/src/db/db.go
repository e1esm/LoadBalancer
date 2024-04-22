package db

import (
	"context"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	hostConf = "hosts.json"
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

	hDB := &HostsDB{
		cli: cli,
	}

	if err := hDB.retry(); err != nil {
		return nil
	}

	return hDB
}

func (h *HostsDB) Set(ctx context.Context, host models.Host) error {
	if host.Address == "" || host.Port == 0 {
		return fmt.Errorf("invalid address")
	}
	return h.cli.Set(ctx, host.Address, host, 0).Err()
}

func (h *HostsDB) GetOrSet(ctx context.Context, addr string) (models.Host, error) {
	var host models.Host
	err := h.cli.Get(ctx, addr).Scan(&host)

	if err != nil {
		addr := strings.Split(addr, ":")
		host.Address, host.Port = addr[0], strToInt(addr[1])
		if err = h.Set(ctx, host); err != nil {
			return models.Host{}, err
		}
	}
	return host, nil
}

func (h *HostsDB) Close() {
	err := h.cli.Close()
	if err != nil {
		log.WithError(err).Error("database connection closing error")
	}
}
