package models

import "encoding/json"

type Host struct {
	Address string `json:"ip"`
	Port    int    `json:"port"`
	Stats   Stats
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
