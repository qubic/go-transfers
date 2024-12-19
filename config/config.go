package config

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

const envPrefix = "QUBIC_TRANSFERS"

type ServerConfig struct {
	HttpHost string
	GrpcHost string
}

type ClientConfig struct {
	EventApiUrl string
	CoreApiUrl  string
}

type DatabaseConfig struct {
	User    string
	Pass    string `conf:"noprint"`
	Host    string
	Port    int
	Name    string
	MaxIdle int
	MaxOpen int
}

type AppConfig struct {
	SyncEnabled bool
	ApiEnabled  bool
}

type LogConfig struct {
	Level     string
	FileError bool
	FileApp   bool
}

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Client   ClientConfig
	Database DatabaseConfig
	Log      LogConfig
}

// GetConfig get config by reading configuration parameters. If available .env files
// will be loaded and used as defaults. See https://github.com/joho/godotenv?tab=readme-ov-file#precedence--conventions
// for used conventions.
func GetConfig(path ...string) (*Config, error) {

	if path == nil || len(path) == 0 {
		loadEnv([]string{""})
	} else {
		loadEnv(path)
	}
	var config Config
	err := conf.Parse(os.Args[1:], envPrefix, &config)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, usageErr := conf.Usage(envPrefix, &config)
			if usageErr != nil {
				slog.Error("generating usage message", "errors", usageErr)
			}
			fmt.Println(usage)
			os.Exit(0)
		}
		if errors.Is(err, conf.ErrVersionWanted) {
			version, versionErr := conf.VersionString(envPrefix, &config)
			if versionErr != nil {
				slog.Error("generating version message", "errors", versionErr)
			}
			fmt.Println(version)
			os.Exit(0)
		}
	}
	return &config, err
}

// loadEnv for conventions see https://github.com/bkeepers/dotenv?tab=readme-ov-file#customizing-rails
func loadEnv(path []string) {
	const ENVIRONMENT = "QUBIC_TRANSFERS_ENV"
	env := os.Getenv(ENVIRONMENT)
	if "" == env {
		env = "development"
	}
	slog.Info("Configuring environment.", "environment", env)

	loadEnvFile(path, ".env."+env+".local")
	if "test" != env {
		// don't use local env overrides for test
		loadEnvFile(path, ".env.local")
	}
	loadEnvFile(path, ".env."+env)
	loadEnvFile(path, ".env")
}

func loadEnvFile(path []string, envFileName string) {
	for _, p := range path {
		realPath := filepath.Join(p, envFileName)
		err := godotenv.Load(realPath)
		if err != nil {
			slog.Debug(fmt.Sprintf("There is no config file [%s].", realPath))
		} else {
			slog.Info("loaded config:", "file", realPath)
		}
	}

}
