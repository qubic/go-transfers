package api

import (
	"flag"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
	"go-transfers/proto"
	"io"
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
	return 1234, nil
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

	if string(body) != "{\"status\":\"UP\"}" {
		t.Errorf("Unexpected response body: [%s]", body)
	}
}

func TestServer_GetAssetTransfersForTick_thenReturnAssetTransfers(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/ticks/1234/events/asset-transfer")
}

func TestServer_GetAssetTransfersForTick_givenUnavailableTickNumber_thenReturnNotFound(t *testing.T) {
	callServiceVerifyStatus(t, "http://localhost:8003/api/v1/ticks/12345/events/asset-transfer", http.StatusNotFound)
}

func TestServer_GetQuTransfersForTick_thenStatusOk(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/ticks/1234/events/qu-transfer")
}

func TestServer_GetQuTransfersForTick_givenUnavailableTickNumber_thenReturnNotFound(t *testing.T) {
	callServiceVerifyStatus(t, "http://localhost:8003/api/v1/ticks/12345/events/qu-transfer", http.StatusNotFound)
}

func TestServer_GetAssetTransfersForEntity_thenReturnAssetTransfers(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/entities/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB/events/asset-transfer")
}

func TestServer_GetAssetTransfersForEntity_GivenInvalidIdentity_thenReturnBadRequest(t *testing.T) {
	callServiceVerifyStatus(t, "http://localhost:8003/api/v1/entities/BLAH/events/asset-transfer", http.StatusBadRequest)
}

func TestServer_GetQuTransfersForEntity_thenStatusOk(t *testing.T) {
	callServiceVerifyNoError(t, "http://localhost:8003/api/v1/entities/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB/events/qu-transfer")
}

func TestServer_GetQuTransfersForEntity_givenInvalidIdentity_thenBadRequest(t *testing.T) {
	callServiceVerifyStatus(t, "http://localhost:8003/api/v1/entities/BLAH/events/qu-transfer", http.StatusBadRequest)
}

func Test_IsValidIdentity(t *testing.T) {
	assert.False(t, isValidIdentity("cfBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL"))
	assert.False(t, isValidIdentity("CFBMEMZOIDEXQAUXYYSZIURADQL123PMNJXQSNVQZAHYVOPYUKKJBJUCTVJL"))
	assert.False(t, isValidIdentity("CFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJLL"))
	assert.False(t, isValidIdentity("CFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJK"))
	assert.False(t, isValidIdentity("AFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL"))
	assert.False(t, isValidIdentity("BMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL"))
	assert.True(t, isValidIdentity("CFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL"))
	assert.True(t, isValidIdentity("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB"))
}

func callServiceVerifyStatus(t *testing.T, url string, expectedStatus int) {
	httpClient := http.DefaultClient
	response, err := httpClient.Get(url)
	if err != nil {
		t.Error(err)
	}
	if response.StatusCode != expectedStatus {
		t.Errorf("Expected status [%d] but got [%s]", expectedStatus, response.Status)
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

func callServiceVerifyNoError(t *testing.T, url string) {
	callServiceVerifyStatus(t, url, http.StatusOK)
}

func readBody(closer io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(closer)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("closing body after reading")
		}
	}(closer)

	return body, err
}
