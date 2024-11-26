package config

import (
	"testing"
)

func Test_GetConfig(t *testing.T) {

	t.Setenv("QUBIC_TRANSFERS_ENV", "test")
	config, getConfigErr := GetConfig()
	if getConfigErr != nil {
		t.Error(getConfigErr)
	}

	var expected = Config{
		Server: ServerConfig{
			HttpHost: "1.2.3.4:5678",
			GrpcHost: "1.2.3.4:6789",
		},
		EventClient: EventClientConfig{
			TargetUrl: "2.3.4.5:6789",
		},
	}

	if *config != expected {
		t.Error("Expected ", expected, "got", config)
	}

}
