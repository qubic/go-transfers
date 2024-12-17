package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-transfers/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
)

type Server struct {
	proto.UnimplementedTransferServiceServer
	listenAddrGRPC string
	listenAddrHTTP string
	repository     Repository
}

type Repository interface {
	GetLatestTick() (int, error)
	GetAssetChangeEventsForTick(tickNumber int) ([]*proto.AssetChangeEvent, error)
	GetQuTransferEventsForTick(tickNumber int) ([]*proto.QuTransferEvent, error)
	GetQuTransferEventsForEntity(identity string) ([]*proto.QuTransferEvent, error)
	GetAssetChangeEventsForEntity(identity string) ([]*proto.AssetChangeEvent, error)
}

func NewServer(grpcAdders, httpAddress string, repository Repository) *Server {

	return &Server{
		listenAddrGRPC: grpcAdders,
		listenAddrHTTP: httpAddress,
		repository:     repository,
	}

}

func (s *Server) Health(_ context.Context, _ *emptypb.Empty) (*proto.HealthResponse, error) {
	return &proto.HealthResponse{
		Status: "UP",
	}, nil // TODO add db connectivity check
}

func (s *Server) GetAssetChangeEventsForTick(_ context.Context, request *proto.TickRequest) (*proto.AssetChangeEventsResponse, error) {
	tickNumber := request.GetTick()
	latestTick, err := s.repository.GetLatestTick() // TODO add context
	// TODO evaluate tick / return 404 if tick is after latest tick
	if err != nil {
		return nil, errorRetrievingEvents("getting latest tick.", "error", err)
	}
	slog.Debug("Get asset transfers:", "tick", tickNumber, "latest", latestTick)
	events, err := s.repository.GetAssetChangeEventsForTick(int(tickNumber)) // TODO add context
	if err != nil {
		return nil, errorRetrievingEvents("getting ownership change events.", "tickNumber", tickNumber, "error", err)
	}
	response := proto.AssetChangeEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetQuTransferEventsForTick(_ context.Context, request *proto.TickRequest) (*proto.QuTransferEventsResponse, error) {
	tickNumber := request.GetTick()
	latestTick, err := s.repository.GetLatestTick() // TODO add context
	// TODO evaluate tick / return 404 if tick is after latest tick
	if err != nil {
		return nil, errorRetrievingEvents("getting latest tick.", "error", err)
	}
	slog.Debug("Get qu transfers:", "tick", tickNumber, "latest", latestTick)
	events, err := s.repository.GetQuTransferEventsForTick(int(tickNumber)) // TODO add context
	if err != nil {
		return nil, errorRetrievingEvents("getting qu transfer events", "tickNumber", tickNumber, "error", err)
	}

	response := proto.QuTransferEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetAssetChangeEventsForEntity(_ context.Context, request *proto.EntityRequest) (*proto.AssetChangeEventsResponse, error) {
	identity := request.GetIdentity()               // FIXME validation
	latestTick, err := s.repository.GetLatestTick() // TODO add context
	if err != nil {
		return nil, errorRetrievingEvents("getting latest tick.", "error", err)
	}
	slog.Debug("Get asset transfers:", "entity", identity, "latest", latestTick)

	events, err := s.repository.GetAssetChangeEventsForEntity(identity)
	if err != nil {
		return nil, errorRetrievingEvents("getting asset change events", "identity", identity, "error", err)
	}

	response := proto.AssetChangeEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetQuTransferEventsForEntity(_ context.Context, request *proto.EntityRequest) (*proto.QuTransferEventsResponse, error) {
	identity := request.GetIdentity()               // FIXME validation
	latestTick, err := s.repository.GetLatestTick() // TODO add context
	if err != nil {
		return nil, errorRetrievingEvents("getting latest tick.", "error", err)
	}
	slog.Debug("Get qu transfers:", "entity", identity, "latest", latestTick)

	events, err := s.repository.GetQuTransferEventsForEntity(identity)
	if err != nil {
		return nil, errorRetrievingEvents("getting qu transfer events", "identity", identity, "error", err)
	}

	response := proto.QuTransferEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func errorRetrievingEvents(internalMessage string, args ...any) error {
	errorId := uuid.New()
	slog.Error(internalMessage, "uuid", errorId.String(), args)
	return status.Errorf(codes.Internal, "error retrieving events [%s]", errorId)
}

func (s *Server) Start() error {
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(600*1024*1024),
		grpc.MaxSendMsgSize(600*1024*1024),
	)
	proto.RegisterTransferServiceServer(srv, s)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", s.listenAddrGRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	if s.listenAddrHTTP != "" {
		go func() {
			mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{EmitDefaultValues: true, EmitUnpopulated: true},
			}))
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(600*1024*1024),
					grpc.MaxCallSendMsgSize(600*1024*1024),
				),
			}

			if err := proto.RegisterTransferServiceHandlerFromEndpoint(
				context.Background(),
				mux,
				s.listenAddrGRPC,
				opts,
			); err != nil {
				panic(err)
			}

			if err := http.ListenAndServe(s.listenAddrHTTP, mux); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}
