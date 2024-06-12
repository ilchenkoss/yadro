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

type UserServer struct {
	auth.UnimplementedUserServer
	us   port.UserService
	slog *slog.Logger
}

func NewUserServer(gRPC *grpc.Server, userService port.UserService, slog *slog.Logger) {
	uServer := &UserServer{
		us:   userService,
		slog: slog,
	}
	auth.RegisterUserServer(gRPC, uServer)
}

func (s *UserServer) UserRole(ctx context.Context, req *auth.UserRoleRequest) (*auth.UserRoleResponse, error) {
	op := "gRPC.Auth.UserRole"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateUserRole(req); vErr != nil {
		return nil, vErr
	}

	userRole, iErr := s.us.UserRole(req.GetUserId())
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

func (s *UserServer) ChangeRole(ctx context.Context, req *auth.ChangeRoleRequest) (*auth.ChangeRoleResponse, error) {
	op := "gRPC.Auth.ChangeRole"
	s.slog.With(slog.String("op", op))

	if vErr := ValidateChangeRole(req); vErr != nil {
		return nil, vErr
	}

	if uurErr := s.us.UpdateUserRole(req.GetReqUserId(), domain.UserRole(req.GetNewUserRole())); uurErr != nil {
		switch {
		case errors.Is(uurErr, domain.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "req user not found")
		case errors.Is(uurErr, domain.ErrUserNotSuperUser):
			return nil, status.Error(codes.PermissionDenied, "change role can only super_user")
		default:
			s.slog.Debug("unhandled error: ", "error", uurErr)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &auth.ChangeRoleResponse{}, nil
}
