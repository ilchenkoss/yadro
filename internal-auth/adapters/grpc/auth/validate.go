package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	auth "myapp/pkg/proto/gen"
)

const (
	emptyValue = 0
)

func ValidateLogin(req *auth.LoginRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func ValidateRegister(req *auth.RegisterRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetRole() == "" {
		return status.Error(codes.InvalidArgument, "role is required")
	}
	return nil
}

func ValidateUserRole(req *auth.UserRoleRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
