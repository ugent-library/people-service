package cli

import (
	"fmt"
	"time"
)

//go:generate go run github.com/g4s8/envdoc@v0.1.2 --output ../CONFIG.md --all

// Version info
type Version struct {
	Branch string `env:"SOURCE_BRANCH"`
	Commit string `env:"SOURCE_COMMIT"`
	Image  string `env:"IMAGE_NAME"`
}

// Application configuration
type Config struct {
	// Env must be local, development, test or production
	Env    string `env:"PEOPLE_ENV" envDefault:"production"`
	Host   string `env:"PEOPLE_HOST"`
	Port   int    `env:"PEOPLE_PORT" envDefault:"3000"`
	APIKey string `env:"PEOPLE_API_KEY,notEmpty"`
	// Repository configuration
	Repo struct {
		// Database connection string
		Conn               string        `env:"CONN,notEmpty"`
		DeactivationPeriod time.Duration `env:"DEACTIVATION_PERIOD,notEmpty" envDefault:"8h"`
	} `envPrefix:"PEOPLE_REPO_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
