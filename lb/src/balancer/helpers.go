package balancer

import (
	"context"
	"github.com/e1esm/LoadBalancer/lb/src/models"
	"strconv"
)

func (lb *LoadBalancer) availableHosts(ctx context.Context) []models.Host {

	hosts := make([]models.Host, 0)

	lb.availableServices.Range(func(key, value any) bool {
		host, err := lb.db.GetOrSet(ctx, key.(string))
		if err != nil {
			return false
		}

		hosts = append(hosts, host)

		return true
	})

	return hosts
}

func strToInt(num string) int {
	n, _ := strconv.Atoi(num)

	return n
}
