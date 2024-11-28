package api

import (
	"flag"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Start server
	srv := NewServer("0.0.0.0:8000", "0.0.0.0:8001")
	err := srv.Start()
	if err != nil {
		os.Exit(-1)
	}

	flag.Parse()
	exitCode := m.Run()

	// Exit
	os.Exit(exitCode)
}

func TestServer_whenHealth_thenReturnStatusUp(t *testing.T) {
	httpClient := http.DefaultClient

	response, err := httpClient.Get("http://localhost:8001/status/health")
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Unexpected response status: [%s]", response.Status)
	}

	body, err := readBody(response.Body)
	if err != nil {
		t.Error(err)
	}

	// TODO how to validate json body?
	if string(body) != "{\"status\":\"UP\"}" {
		t.Errorf("Unexpected response body: [%s]", body)
	}
}

func TestServer_GetAssetTransfersForTick_thenReturnAssetTransfers(t *testing.T) {
	httpClient := http.DefaultClient
	response, err := httpClient.Get("http://localhost:8001/v1/tick/1234/asset-transfers")
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Unexpected response status: [%s]", response.Status)
	}

	body, err := readBody(response.Body)
	if err != nil {
		t.Error(err)
	}

	if body == nil {
		t.Error("Response body is nil")
	}
}

func readBody(closer io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(closer)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO log here?
		}
	}(closer)

	return body, err
}
