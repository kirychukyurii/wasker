package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	userSvc "github.com/kirychukyurii/wasker/gen/go/user/v1alpha1"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type HttpServer struct {
	Server *http.Server
}

func NewHttpServer(logger log.Logger) HttpServer {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux()

	// setting up a dail up for gRPC service by specifying endpoint/target url
	err := userSvc.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, "localhost:8080", opts)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error registering handler from endpoint")
	}

	// Creating a normal HTTP server
	server := &http.Server{
		Handler: mux,
	}

	return HttpServer{
		Server: server,
	}
}
