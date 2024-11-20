package client

import (
	"context"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func logEvents() error {

	// post tick number in body
	conn, err := grpc.NewClient("95.216.243.140:8003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.Wrap(err, "creating grpc connection")
	}

	eventsClient := eventspb.NewEventsServiceClient(conn)
	tickEvents, err := eventsClient.GetTickEvents(context.Background(), &eventspb.GetTickEventsRequest{Tick: 17275986})
	if err != nil {
		return errors.Wrap(err, "getting tick events")
	}

	log.Printf("%+v\n", tickEvents)

	return nil

}
