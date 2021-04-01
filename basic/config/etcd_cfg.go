package config

type EtcdConfig interface {
	GetEnabled() bool
	GetPort() int
	GetHost() string
}

type defaultEtcdConfig struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

func (c defaultEtcdConfig) GetPort() int {
	return c.Port
}

func (c defaultEtcdConfig) GetEnabled() bool {
	return c.Enabled
}

func (c defaultEtcdConfig) GetHost() string {
	return c.Host
}
