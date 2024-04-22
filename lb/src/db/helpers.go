package db

import (
	"context"
	"fmt"
	"strconv"
)

const maxAttempts = 10

func (h *HostsDB) retry() error {
	var attempt int

	for attempt < maxAttempts {
		if status := h.cli.Ping(context.Background()); status.Err() == nil {
			return nil
		}
		attempt++
	}

	return fmt.Errorf("no connection was established")
}

func strToInt(num string) int {
	n, _ := strconv.Atoi(num)

	return n
}
