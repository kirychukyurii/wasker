package server

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/wasker/internal/constants"
	"github.com/kirychukyurii/wasker/internal/gateway/server/register"
	"github.com/kirychukyurii/wasker/internal/pkg/consul"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type HttpServer struct {
	Server *http.Server
}

func NewHttpServer(logger log.Logger, discovery consul.ServiceDiscovery) HttpServer {
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	services, err := discovery.GetByName(constants.DirectoryServiceName)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("get services from discovery")
	}

	directoryService := services[0]

	// setting up a dail up for gRPC service by specifying endpoint/target url
	if err := register.GrpcDirectoryEndpoints(ctx, mux, fmt.Sprintf("%s:%d", directoryService.Host, directoryService.Port), opts); err != nil {
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
