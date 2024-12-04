package config

import (
	"log/slog"
	"testing"
)

func TestConfig_GetConfig(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

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
	}

	if *config != expected {
		t.Error("Expected ", expected, "got", config)
	}

}
