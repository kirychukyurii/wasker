package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/server/register"
)

type HttpServer struct {
	Server *http.Server
}

func NewHttpServer(logger log.Logger) HttpServer {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()

	// setting up a dail up for gRPC service by specifying endpoint/target url
	if err := register.GrpcDirectoryEndpoints(ctx, mux, "localhost:8080", opts); err != nil {
		logger.Log.Fatal().Err(err).Msg("Error registering handlers from directory endpoint")
	}

	// Creating a normal HTTP server
	server := &http.Server{
		Handler: mux,
	}

	return HttpServer{
		Server: server,
	}
}
