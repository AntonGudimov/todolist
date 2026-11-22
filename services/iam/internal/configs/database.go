package configs

import "fmt"

type DatabaseConfig struct {
	DatabaseUser     string `env:"IAM_DATABASE_USER"`
	DatabasePassword string `env:"IAM_DATABASE_PASSWORD"`
	DatabaseHost     string `env:"IAM_DATABASE_HOST"`
	DatabasePort     int    `env:"IAM_DATABASE_PORT"`
	DatabaseName     string `env:"IAM_DATABASE_NAME"`
}

func (c *DatabaseConfig) GetURL() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}
