package cli

import "fmt"

// Version info
type Version struct {
	Branch string `env:"SOURCE_BRANCH"`
	Commit string `env:"SOURCE_COMMIT"`
	Image  string `env:"IMAGE_NAME"`
}

type Config struct {
	// Env must be local, development, test or production
	Env    string `env:"PEOPLE_ENV" envDefault:"production"`
	Host   string `env:"PEOPLE_HOST"`
	Port   int    `env:"PEOPLE_PORT" envDefault:"3000"`
	APIKey string `env:"PEOPLE_API_KEY,notEmpty"`
	Repo   struct {
		Conn string `env:"CONN,notEmpty"`
	} `envPrefix:"PEOPLE_REPO_"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
