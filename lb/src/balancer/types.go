package balancer

import "fmt"

type Host struct {
	ip   string
	port int
}

func (h *Host) Address() string {
	return fmt.Sprintf("%s:%d", h.ip, h.port)
}
