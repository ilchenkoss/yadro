package app

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"myapp/internal-auth/adapters/database"
	"myapp/internal-auth/adapters/database/repository"
	"myapp/internal-auth/adapters/grpc/auth"
	"myapp/internal-auth/config"
	"myapp/internal-auth/core/service"
	"net"
)

type App struct {
	slog       *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	slog *slog.Logger, cfg config.Config) *App {
	gRPCServer := grpc.NewServer()
	dbConnection, cErr := database.NewConnection(&cfg)
	if cErr != nil {
		panic(cErr)
	}
	pErr := dbConnection.Ping()
	if pErr != nil {
		panic(pErr)
	}
	mErr := dbConnection.MakeMigrations()
	if mErr != nil {
		panic(mErr)
	}

	userRepo := repository.NewUserRepository(dbConnection)

	tokenService := service.NewTokenService(cfg)
	authService := service.NewAuthService(userRepo, tokenService)

	auth.Register(gRPCServer, authService, slog)

	return &App{
		slog:       slog,
		gRPCServer: gRPCServer,
		port:       cfg.Server.Port,
	}
}

func (a *App) AppRun() error {
	return a.GRPCRun()
}

func (a *App) GRPCRun() error {
	op := "app.GRPCRun"
	a.slog.With(slog.String("op", op))

	listener, lErr := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))
	if lErr != nil {
		a.slog.Error("error init listener")
		return fmt.Errorf("%s: %w", op, lErr)
	}
	slog.Info("grpc server is listening", slog.String("addr", listener.Addr().String()))

	if sErr := a.gRPCServer.Serve(listener); sErr != nil {
		return fmt.Errorf("%s: %w", op, sErr)
	}

	return nil
}

func (a *App) GRPCStop() {
	op := "app.GRPCStop"
	a.slog.With(slog.String("op", op)).
		Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
