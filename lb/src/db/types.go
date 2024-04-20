package db

type configuredHosts struct {
	Servers []server `json:"servers"`
}

type server struct {
	Address string `json:"ip"`
	Port    int    `json:"port"`
}
