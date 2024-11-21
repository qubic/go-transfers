package client

import (
	"context"
	eventspb "github.com/qubic/go-events/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventClient struct {
	protoClient eventspb.EventsServiceClient
}

func NewEventClient(targetUrl string) (*EventClient, error) {
	conn, err := grpc.NewClient(targetUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	e := &EventClient{
		protoClient: eventspb.NewEventsServiceClient(conn),
	}
	return e, err
}

func (eventClient *EventClient) GetEvents(tickNumber uint32) (*eventspb.TickEvents, error) {
	return eventClient.protoClient.GetTickEvents(context.Background(), &eventspb.GetTickEventsRequest{Tick: tickNumber})
}
