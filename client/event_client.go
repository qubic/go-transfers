package client

import (
	"context"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	qubicpb "github.com/qubic/go-qubic/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IntegrationEventClient struct {
	eventApi eventspb.EventsServiceClient
	coreApi  qubicpb.CoreServiceClient
}

type TickInfo struct {
	CurrentTick uint32
	InitialTick uint32
}

type EventStatus struct {
	AvailableTick uint32
}

func NewIntegrationEventClient(eventApiUrl, coreApiUrl string) (*IntegrationEventClient, error) {
	eventApiConn, err := grpc.NewClient(eventApiUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "creating event api connection")
	}
	coreApiConn, err := grpc.NewClient(coreApiUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "creating core api connection")
	}
	e := IntegrationEventClient{
		eventApi: eventspb.NewEventsServiceClient(eventApiConn),
		coreApi:  qubicpb.NewCoreServiceClient(coreApiConn),
	}
	return &e, nil
}

func (eventClient *IntegrationEventClient) GetEvents(context context.Context, tickNumber uint32) (*eventspb.TickEvents, error) {
	return eventClient.eventApi.GetTickEvents(context, &eventspb.GetTickEventsRequest{Tick: tickNumber})
}

func (eventClient *IntegrationEventClient) GetStatus(context context.Context) (*EventStatus, error) {
	s, err := eventClient.eventApi.GetStatus(context, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting events status")
	}
	status := EventStatus{
		AvailableTick: s.GetLastProcessedTick().GetTickNumber(),
	}
	return &status, nil
}

func (eventClient *IntegrationEventClient) GetTickInfo(context context.Context) (*TickInfo, error) {
	ti, err := eventClient.coreApi.GetTickInfo(context, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting tick info")
	}
	tiDto := TickInfo{
		CurrentTick: ti.Tick,
		InitialTick: ti.InitialTickOfEpoch,
	}
	return &tiDto, nil
}
