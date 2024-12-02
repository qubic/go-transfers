package client

import (
	"context"
	eventspb "github.com/qubic/go-events/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IntegrationEventClient struct {
	protoClient eventspb.EventsServiceClient
}

func NewIntegrationEventClient(targetUrl string) (*IntegrationEventClient, error) {
	conn, err := grpc.NewClient(targetUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	e := &IntegrationEventClient{
		protoClient: eventspb.NewEventsServiceClient(conn),
	}
	return e, err
}

func (eventClient *IntegrationEventClient) GetEvents(tickNumber uint32) (*eventspb.TickEvents, error) {
	return eventClient.protoClient.GetTickEvents(context.Background(), &eventspb.GetTickEventsRequest{Tick: tickNumber})
}

func (eventClient *IntegrationEventClient) GetStatus() (*eventspb.GetStatusResponse, error) {
	return eventClient.protoClient.GetStatus(context.Background(), &emptypb.Empty{})
}
