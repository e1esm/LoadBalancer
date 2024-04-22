package db

import (
	"context"
	"fmt"
)

const maxAttempts = 10

func (h *HostsDB) retry() error {
	var attempt int

	for attempt < maxAttempts {
		if err := h.cli.Ping(context.Background()); err == nil {
			return nil
		}
		attempt++
	}

	return fmt.Errorf("no connection was established")
}
