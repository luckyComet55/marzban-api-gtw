package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	srv "github.com/luckyComet55/marzban-api-gtw/internal/gateway_server"
	pcl "github.com/luckyComet55/marzban-api-gtw/internal/panel_client"
)

type Application struct {
	logger     *slog.Logger
	grpcServer *grpc.Server
	port       uint64
}

func NewApp(
	serverLogger *slog.Logger,
	clientLogger *slog.Logger,
	username string,
	password string,
	marzbanBaseUrl string,
	port uint64,
) *Application {
	cli := pcl.NewMarzbanPanelClient(pcl.MarzbanPanelClientConfig{
		Username:       username,
		Password:       password,
		MarzbanBaseUrl: marzbanBaseUrl,
	}, clientLogger)
	grpcServer := grpc.NewServer()
	srv.Register(grpcServer, cli, serverLogger)
	reflection.Register(grpcServer)
	return &Application{
		logger:     serverLogger,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (app *Application) MustStart() {
	if err := app.Start(); err != nil {
		panic(err)
	}
}

func (app *Application) Start() error {
	const op string = "grcpapp.Start"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	app.logger.With(slog.String("op", op)).
		Info("gRPC server started", slog.String("addr", l.Addr().String()))

	if err := app.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *Application) Stop() {
	const op string = "grpcapp.Stop"

	app.logger.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Uint64("port", app.port))

	app.grpcServer.GracefulStop()
}
