package config

import (
	"github.com/gookit/slog"
	"testing"
)

func TestConfig_GetConfig(t *testing.T) {
	slog.SetLogLevel(slog.DebugLevel)

	t.Setenv("QUBIC_TRANSFERS_ENV", "test")
	config, getConfigErr := GetConfig(".", "..")
	if getConfigErr != nil {
		t.Error(getConfigErr)
	}

	var expected = Config{
		Server: ServerConfig{
			HttpHost: "1.2.3.4:5678", // .env.test
			GrpcHost: "1.2.3.4:6789", // .env.test
		},
		Client: ClientConfig{
			EventApiUrl: "2.3.4.5:6789", // .env.test
			CoreApiUrl:  "2.3.4.5:5678",
		},
		Database: DatabaseConfig{
			User:    "test",      // .env.test
			Pass:    "test-pass", // .env.test
			Host:    "localhost", // global default
			Port:    5432,        // global default
			Name:    "test",      // .env.test
			MaxIdle: 10,
			MaxOpen: 10,
		},
		App: AppConfig{
			SyncEnabled: true,
			ApiEnabled:  true,
		},
		Log: LogConfig{
			Level:     "Info",
			FileApp:   false,
			FileError: false,
		},
	}

	if *config != expected {
		t.Error("Expected ", expected, "got", config)
	}

}
