package balancer

import "fmt"

type Host struct {
	ip   string
	port int
	n    int
}

func (h *Host) Address() string {
	return fmt.Sprintf("%s-%d:%d", h.ip, h.n, h.port)
}
