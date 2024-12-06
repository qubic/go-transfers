package api

import (
	"flag"
	"go-transfers/proto"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
)

type FakeRepository struct {
}

func (f FakeRepository) GetAssetChangeEventsForEntity(_ string) ([]*proto.AssetChangeEvent, error) {
	return []*proto.AssetChangeEvent{}, nil
}

func (f FakeRepository) GetQuTransferEventsForEntity(_ string) ([]*proto.QuTransferEvent, error) {
	return []*proto.QuTransferEvent{}, nil
}

func (f FakeRepository) GetAssetChangeEventsForTick(_ int) ([]*proto.AssetChangeEvent, error) {
	return []*proto.AssetChangeEvent{}, nil
}

func (f FakeRepository) GetQuTransferEventsForTick(_ int) ([]*proto.QuTransferEvent, error) {
	return []*proto.QuTransferEvent{}, nil
}

func (f FakeRepository) GetLatestTick() (int, error) {
	return 123, nil
}

func TestMain(m *testing.M) {

	// Start server
	srv := NewServer("0.0.0.0:8002", "0.0.0.0:8003", &FakeRepository{})
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

	response, err := httpClient.Get("http://localhost:8003/status/health")
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
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/ticks/1234/events/asset-transfer")
}

func TestServer_GetQuTransfersForTick_thenStatusOk(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/ticks/1234/events/qu-transfer")
}

func TestServer_GetAssetTransfersForEntity_thenReturnAssetTransfers(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/entities/BLAH/events/asset-transfer")
}

func TestServer_GetQuTransfersForEntity_thenStatusOk(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/entities/BLAH/events/qu-transfer")
}

func callServiceVerifyNoError(t *testing.T, url string) {
	httpClient := http.DefaultClient
	response, err := httpClient.Get(url)
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
	slog.Info("Read body.", "body", body)
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
