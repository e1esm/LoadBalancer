package models

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Address string `json:"ip"`
	Port    int    `json:"port"`
	Stats   Stats
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%d", h.Address, h.Port)
}

func (h Host) MarshalBinary() ([]byte, error) {
	return json.Marshal(h)
}

func (h *Host) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, h)
}

type Stats struct {
	OngoingReqs     int
	OverallRequests int
}
