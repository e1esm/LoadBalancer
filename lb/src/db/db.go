package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
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

	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.WithError(err).Fatal("could not connect to redis")
	}

	hdb := &HostsDB{
		cli: cli,
	}

	if err := hdb.updateConf(hostConf); err != nil {
		log.WithError(err).Fatal("configuration save error")
	}

	return hdb
}

func (h *HostsDB) GetHosts(ctx context.Context) ([]models.Host, error) {
	cmd := h.cli.Do(ctx, "KEYS", "*")

	res, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("error while fetching values: %w", err)
	}

	hosts := make([]models.Host, 0)

	for _, entity := range res.([]interface{}) {
		host, err := h.Get(context.Background(), entity.(string))
		if err != nil {
			return nil, fmt.Errorf("error while getting host entity: %w", err)
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}

func (h *HostsDB) Set(ctx context.Context, host models.Host) error {
	return h.cli.Set(ctx, host.Address, host, 0).Err()
}

func (h *HostsDB) Get(ctx context.Context, addr string) (models.Host, error) {
	var host models.Host
	err := h.cli.Get(ctx, addr).Scan(&host)
	if err != nil {
		return models.Host{}, fmt.Errorf("scanning error: %w", err)
	}
	return host, nil
}

func (h *HostsDB) Close() {
	err := h.cli.Close()
	if err != nil {
		log.WithError(err).Error("database connection closing error")
	}
}

func (h *HostsDB) updateConf(confFile string) error {
	f, err := os.Open(confFile)
	if err != nil {
		return fmt.Errorf("updating conf error: %w", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading conf file error: %w", err)
	}

	var availableHosts configuredHosts

	if err := json.Unmarshal(b, &availableHosts); err != nil {
		return fmt.Errorf("marshalling error: %w", err)
	}

	for _, host := range availableHosts.Servers {

		dbH := models.Host{
			Address: host.Address,
			Port:    host.Port,
		}

		log.Infof("host: %s:%d", dbH.Address, dbH.Port)

		if err := h.Set(context.Background(), dbH); err != nil {
			return fmt.Errorf("saving host error: %w", err)
		}
	}

	return nil
}
