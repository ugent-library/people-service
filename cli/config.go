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
	Env    string `env:"ENV" envDefault:"production"`
	Host   string `env:"HOST"`
	Port   int    `env:"PORT" envDefault:"3000"`
	APIKey string `env:"API_KEY,notEmpty"`
	// Repository configuration
	Repo struct {
		// Database connection string
		Conn               string        `env:"CONN,notEmpty"`
		DeactivationPeriod time.Duration `env:"DEACTIVATION_PERIOD" envDefault:"8h"`
	} `envPrefix:"REPO_"`
	// Search index configuration
	Index struct {
		// Connection string
		Conn string `env:"CONN,notEmpty"`
		// Index Name
		Name string `env:"NAME,notEmpty"`
		// Index retention
		Retention int `env:"RETENTION" envDefault:"5"`
	} `envPrefix:"INDEX_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
