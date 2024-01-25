package cli

import "fmt"

type ConfigDb struct {
	Url    string `env:"URL" envDefault:"postgres://people:people@localhost:5432/authority?sslmode=disable"`
	AesKey string `env:"AES_KEY,notEmpty"`
}

type ConfigApi struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port int    `env:"PORT" envDefault:"3999"`
	Key  string `env:"KEY,notEmpty"`
}

type ConfigLdap struct {
	Url      string `env:"URL,notEmpty"`
	Username string `env:"USERNAME,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
}

type Config struct {
	Production bool       `env:"PRODUCTION"`
	Db         ConfigDb   `envPrefix:"DB_"`
	Api        ConfigApi  `envPrefix:"API_"`
	Ldap       ConfigLdap `envPrefix:"LDAP_"`
	IPRanges   string     `env:"IP_RANGES"`
}

func (ca ConfigApi) Addr() string {
	return fmt.Sprintf("%s:%d", ca.Host, ca.Port)
}
