package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"myapp/internal-auth/core/domain"
	"myapp/internal-auth/core/port"
	auth "myapp/pkg/proto/gen"
)

type server struct {
	auth.UnimplementedAuthServer
	as   port.AuthService
	slog *slog.Logger
}

func Register(gRPC *grpc.Server, authService port.AuthService, slog *slog.Logger) {
	auth.RegisterAuthServer(gRPC, &server{
		as:   authService,
		slog: slog})
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	op := "gRPC.Auth.Login"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateLogin(req); vErr != nil {
		return nil, vErr
	}

	token, lErr := s.as.Login(req.GetLogin(), req.GetPassword())
	if lErr != nil {
		switch {
		case errors.Is(lErr, domain.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			s.slog.Debug("unhandled error: ", "error", lErr)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}

func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	op := "gRPC.Auth.Register"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateRegister(req); vErr != nil {
		return nil, vErr
	}

	userID, rErr := s.as.Register(req.GetLogin(), req.GetPassword(), domain.UserRole(req.GetRole()))
	if rErr != nil {
		switch {
		case errors.Is(rErr, domain.ErrUserAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		default:
			s.slog.Debug("unhandled error: ", "error", rErr)
			return nil, status.Error(codes.Internal, "internal error")
		}

	}
	return &auth.RegisterResponse{UserId: userID}, nil
}

func (s *server) UserRole(ctx context.Context, req *auth.UserRoleRequest) (*auth.UserRoleResponse, error) {
	op := "gRPC.Auth.UserRole"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateUserRole(req); vErr != nil {
		return nil, vErr
	}

	userRole, iErr := s.as.UserRole(req.GetUserId())
	if iErr != nil {
		switch {
		case errors.Is(iErr, domain.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			s.slog.Debug("unhandled error: ", "error", iErr)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &auth.UserRoleResponse{UserRole: string(userRole)}, nil
}

func (s *server) UserID(ctx context.Context, req *auth.UserIDRequest) (*auth.UserIDResponse, error) {
	op := "gRPC.Auth.UserID"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateUserID(req); vErr != nil {
		return nil, vErr
	}

	userID, uIDErr := s.as.UserID(req.GetToken())
	if uIDErr != nil {
		switch {
		case errors.Is(uIDErr, domain.ErrTokenExpired):
			return nil, status.Error(codes.DeadlineExceeded, "token expired")
		case errors.Is(uIDErr, domain.ErrTokenNotValid):
			return nil, status.Error(codes.Unauthenticated, "token not valid")
		default:
			s.slog.Debug("unhandled error: ", "error", uIDErr)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &auth.UserIDResponse{UserId: userID}, nil
}
