package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	ownerrpc "github.com/viacheslavek/grpcauth/auth/internal/grpc/ownerCtl"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, ownerService ownerrpc.OwnerCtl, port int) *App {
	gRPCServer := grpc.NewServer()

	ownerrpc.Register(gRPCServer, ownerService, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run app server", err.Error())
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("starting gRPC server")

	l, errL := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if errL != nil {
		return fmt.Errorf("%s: %w", op, errL)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s failed serve grpc: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
