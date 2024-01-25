package cli

import (
	"github.com/caarlos0/env/v8"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	// load .env file if present
	_ "github.com/joho/godotenv/autoload"
)

var (
	config  Config
	version Version
)

var logger *zap.SugaredLogger

var rootCmd = &cobra.Command{
	Use: "people-service",
}

func init() {
	cobra.OnInitialize(initVersion, initConfig, initLogger)
	cobra.OnFinalize(func() {
		logger.Sync()
	})
}

func initConfig() {
	cobra.CheckErr(env.ParseWithOptions(&config, env.Options{
		Prefix: "PEOPLE_",
	}))
}

func initVersion() {
	cobra.CheckErr(env.Parse(&version))
}

func initLogger() {
	var l *zap.Logger
	var e error
	if config.Production {
		l, e = zap.NewProduction()
	} else {
		l, e = zap.NewDevelopment()
	}
	cobra.CheckErr(e)
	logger = l.Sugar()
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
