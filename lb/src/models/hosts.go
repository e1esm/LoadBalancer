package models

type Host struct {
	Address string
	Port    int
	Stats   Stats
}

type Stats struct {
	OngoingReqs     int
	OverallRequests int
}
