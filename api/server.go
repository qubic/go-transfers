package api

import (
	"context"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/qubic/go-qubic/common"
	"go-transfers/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	proto.UnimplementedTransferServiceServer
	listenAddrGRPC string
	listenAddrHTTP string
	repository     Repository
}

type Repository interface {
	GetLatestTick(ctx context.Context) (int, error)
	GetAssetChangeEventsForTick(ctx context.Context, tickNumber int) ([]*proto.AssetChangeEvent, error)
	GetQuTransferEventsForTick(ctx context.Context, tickNumber int) ([]*proto.QuTransferEvent, error)
	GetQuTransferEventsForEntity(ctx context.Context, identity string) ([]*proto.QuTransferEvent, error)
	GetAssetChangeEventsForEntity(ctx context.Context, identity string) ([]*proto.AssetChangeEvent, error)
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

func (s *Server) GetAssetChangeEventsForTick(ctx context.Context, request *proto.TickRequest) (*proto.AssetChangeEventsResponse, error) {
	tickNumber := request.GetTick()
	latestTick, err := s.repository.GetLatestTick(ctx)
	if latestTick < int(tickNumber) {
		return nil, tickNotFound(tickNumber, latestTick)
	}
	if err != nil {
		return nil, retrieveEventsError("getting latest tick.", "error", err)
	}
	slog.Debug("Get asset transfers:", "tick", tickNumber, "latest", latestTick)
	events, err := s.repository.GetAssetChangeEventsForTick(ctx, int(tickNumber))
	if err != nil {
		return nil, retrieveEventsError("getting ownership change events.", "tickNumber", tickNumber, "error", err)
	}
	response := proto.AssetChangeEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetQuTransferEventsForTick(ctx context.Context, request *proto.TickRequest) (*proto.QuTransferEventsResponse, error) {
	tickNumber := request.GetTick()
	latestTick, err := s.repository.GetLatestTick(ctx)
	if latestTick < int(tickNumber) {
		return nil, tickNotFound(tickNumber, latestTick)
	}
	if err != nil {
		return nil, retrieveEventsError("getting latest tick.", "error", err)
	}
	slog.Debug("Get qu transfers:", "tick", tickNumber, "latest", latestTick)
	events, err := s.repository.GetQuTransferEventsForTick(ctx, int(tickNumber))
	if err != nil {
		return nil, retrieveEventsError("getting qu transfer events", "tickNumber", tickNumber, "error", err)
	}

	response := proto.QuTransferEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetAssetChangeEventsForEntity(ctx context.Context, request *proto.EntityRequest) (*proto.AssetChangeEventsResponse, error) {
	identity := request.GetIdentity()
	if !isValidIdentity(identity) {
		return nil, invalidIdentity(identity)
	}
	latestTick, err := s.repository.GetLatestTick(ctx)
	if err != nil {
		return nil, retrieveEventsError("getting latest tick.", "error", err)
	}
	slog.Debug("Get asset transfers:", "entity", identity, "latest", latestTick)

	events, err := s.repository.GetAssetChangeEventsForEntity(ctx, identity)
	if err != nil {
		return nil, retrieveEventsError("getting asset change events", "identity", identity, "error", err)
	}

	response := proto.AssetChangeEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func (s *Server) GetQuTransferEventsForEntity(ctx context.Context, request *proto.EntityRequest) (*proto.QuTransferEventsResponse, error) {
	identity := request.GetIdentity()
	if !isValidIdentity(identity) {
		return nil, invalidIdentity(identity)
	}
	latestTick, err := s.repository.GetLatestTick(ctx)
	if err != nil {
		return nil, retrieveEventsError("getting latest tick.", "error", err)
	}
	slog.Debug("Get qu transfers", "entity", identity, "latest", latestTick)

	events, err := s.repository.GetQuTransferEventsForEntity(ctx, identity)
	if err != nil {
		return nil, retrieveEventsError("getting qu transfer events", "identity", identity, "error", err)
	}

	response := proto.QuTransferEventsResponse{LatestTick: uint32(latestTick), Events: events}
	return &response, nil
}

func isValidIdentity(s string) bool {
	if len(s) == 60 && !strings.ContainsFunc(s, func(r rune) bool {
		return r < 'A' || r > 'Z'
	}) {
		id := common.Identity(s)
		pubKey, err := id.ToPubKey(false)
		if err != nil {
			return false
		}
		err = id.FromPubKey(pubKey, false)
		return err == nil && id.String() == s
	}
	return false
}

func invalidIdentity(id string) error {
	errorId := uuid.New().String()
	slog.Error("invalid request", "identity", id, "uuid", errorId)
	return status.Errorf(codes.InvalidArgument, "invalid identity [%s]", errorId)
}

func tickNotFound(requested uint32, latestAvailable int) error {
	errorId := uuid.New().String()
	slog.Error("tick not found.", "requested:", requested, "latest:", latestAvailable, "uuid:", errorId)
	return status.Errorf(codes.NotFound, "tick not found. [%s]", errorId)
}

func retrieveEventsError(internalMessage string, args ...any) error {
	errorId := uuid.New().String()
	slog.Error(internalMessage, "uuid", errorId, args)
	return status.Errorf(codes.Internal, "error retrieving events. [%s]", errorId)
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
		slog.Fatalf("failed to listen: %v", err)
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
