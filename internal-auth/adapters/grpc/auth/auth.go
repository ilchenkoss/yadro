package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	auth "myapp/internal-auth/adapters/grpc/proto/gen"
	"myapp/internal-auth/core/domain"
	"myapp/internal-auth/core/port"
	auth "myapp/pkg/proto/gen"
)

type server struct {
	auth.UnimplementedAuthServer
	as port.AuthService
}

func Register(gRPC *grpc.Server, authService port.AuthService) {
	auth.RegisterAuthServer(gRPC, &server{
		as: authService})
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	vErr := ValidateLogin(req)
	if vErr != nil {
		return nil, vErr
	}

	token, lErr := s.as.Login(req.GetLogin(), req.GetPassword())
	if lErr != nil {
		switch {
		case errors.Is(lErr, domain.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}

func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	vErr := ValidateRegister(req)
	if vErr != nil {
		return nil, vErr
	}
	userID, rErr := s.as.Register(req.GetLogin(), req.GetPassword(), domain.UserRole(req.GetRole()))
	if rErr != nil {
		switch {
		case errors.Is(rErr, domain.ErrUserAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}

	}
	return &auth.RegisterResponse{UserId: userID}, nil
}

func (s *server) UserRole(ctx context.Context, req *auth.UserRoleRequest) (*auth.UserRoleResponse, error) {
	vErr := ValidateUserRole(req)
	if vErr != nil {
		return nil, vErr
	}
	userRole, iErr := s.as.UserRole(req.GetUserId())
	if iErr != nil {
		switch {
		case errors.Is(iErr, domain.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &auth.UserRoleResponse{UserRole: string(userRole)}, nil
}
