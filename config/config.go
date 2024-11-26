package config

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

const envPrefix = "QUBIC_TRANSFERS"

type ServerConfig struct {
	HttpHost string
	GrpcHost string
}

type EventClientConfig struct {
	TargetUrl string
}

type DatabaseConfig struct {
	User string
	Pass string `conf:"noprint"`
	Url  string
}

type Config struct {
	Server      ServerConfig
	EventClient EventClientConfig
	Database    DatabaseConfig
}

var lock = &sync.Mutex{}
var loadedConfig *Config = nil

func GetConfig() (*Config, error) {
	if loadedConfig == nil {
		lock.Lock()
		defer lock.Unlock()
		if loadedConfig == nil {
			loadEnv()
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
				return &config, err
			}
			loadedConfig = &config
		}
	}
	return loadedConfig, nil
}

// loadEnv for conventions see https://github.com/bkeepers/dotenv?tab=readme-ov-file#customizing-rails
func loadEnv() {
	const ENVIRONMENT = "QUBIC_TRANSFERS_ENV"
	env := os.Getenv(ENVIRONMENT)
	if "" == env {
		env = "development"
	}
	slog.Info("Configuring environment.", "environment", env)

	loadEnvFile(".env." + env + ".local")
	if "test" != env {
		// don't use local env overrides for test
		loadEnvFile(".env.local")
	}
	loadEnvFile(".env." + env)
	loadEnvFile(".env")
}

func loadEnvFile(path string) {
	realPath := getEnvPath(path)
	err := godotenv.Load(realPath)
	if err != nil {
		slog.Debug(fmt.Sprintf("There is no config file [%s].", realPath))
	} else {
		slog.Info("Loading config.", "file", realPath)
	}
}

// getEnvPath returns the absolute path of the given environment file (envFile).
// It searches for either the requested envFile or the 'go.mod' file from the current working directory upwards
// Assumption is that env file needs at least be in the root ('go.mod') directory, if there is such an env file.
func getEnvPath(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get current directory", "error", err)
		panic(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, envFile)); err == nil {
			break // OK found
		}
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// reached go.mod dir. Don't go further.
			slog.Debug("Go mod dir reached.", "path", goModPath)
			break
		}
		parent := filepath.Dir(currentDir)
		if string(os.PathSeparator) == parent {
			// reached root dir. Can't go further
			slog.Debug("Root dir reached.", "file", envFile)
			break
		} else {
			currentDir = parent
		}
	}

	return filepath.Join(currentDir, envFile)
}
