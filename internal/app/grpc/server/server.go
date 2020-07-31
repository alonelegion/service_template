package server

import (
	"context"
	"github.com/alonelegion/service_template/internal/app/config"
	"github.com/chapsuk/wait"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	servicetemplateapi "github.com/alonelegion/service_template/api"
	"go.elastic.co/apm/module/apmgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"strconv"
)

type (
	GRPC struct {
		logger *zap.Logger
		config *config.AppConfig
	}
)

func NewServer(
	ctx context.Context,
	logger *zap.Logger,
	config *config.AppConfig,
) <-chan error {
	return Register(func() error {
		return (&GRPC{
			logger: logger,
			config: config,
		}).Start(ctx)
	})
}

func (g *GRPC) Start(ctx context.Context) error {
	address := ":" + strconv.Itoa(g.config.GRPC.Port)
	conn, err := net.Listen("tcp4", address)
	if err != nil {
		g.logger.Fatal("error while listen socket for grpc service", zap.Error(err))
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_prometheus.UnaryServerInterceptor,
				apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			),
		),
	)

	servicetemplateapi.RegisterHelloWorldServer(server, g)
	healthpb.RegisterHealthServer(server, health.NewServer())

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(server)

	reflection.Register(server)

	g.logger.Info("Start grpc server", zap.String("address", address))
	select {
	case err := <-Register(func() error {
		return server.Serve(conn)
	}):
		return err
	case <-ctx.Done():
		g.logger.Info("Shutdown grpc server")
		server.GracefulStop()

		return ctx.Err()
	}
}

func Register(fn func() error) <-chan error {
	ch := make(chan error)
	wg := wait.Group{}
	wg.Add(func() {
		ch <- fn()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func (g *GRPC) Get(ctx context.Context, request *servicetemplateapi.Request) (*servicetemplateapi.Response, error) {
	return &servicetemplateapi.Response{Message: "Hello World!"}, nil
}
