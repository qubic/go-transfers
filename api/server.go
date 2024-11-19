package api

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-transfers/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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
}

func NewServer(grpcAdders, httpAddress string) *Server {

	return &Server{
		listenAddrGRPC: grpcAdders,
		listenAddrHTTP: httpAddress,
	}

}

func (s *Server) Health(_ context.Context, _ *emptypb.Empty) (*proto.HealthResponse, error) {
	return &proto.HealthResponse{
		Status: "UP",
	}, nil
}

func (s *Server) GetAssetTransfersForTick(_ context.Context, request *proto.TickRequest) (*proto.AssetTransferResponse, error) {
	tickNumber := request.GetTick()
	log.Printf("tick number %d", tickNumber)
	return &proto.AssetTransferResponse{}, nil
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
