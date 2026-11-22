package configs

import "fmt"

type HostConfig struct {
	HostName string `env:"HOST_NAME"`
	HostPort int    `env:"HOST_PORT"`
}

func (c *HostConfig) GetUrl() string {
	return fmt.Sprintf(
		"%v:%v",
		c.HostName,
		c.HostPort,
	)
}
