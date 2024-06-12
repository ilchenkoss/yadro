package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"myapp/internal-auth/core/domain"
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

func ValidateUserID(req *auth.UserIDRequest) error {
	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, "token is required")
	}
	return nil
}

func ValidateChangeRole(req *auth.ChangeRoleRequest) error {
	if req.GetNewUserRole() == "" {
		return status.Error(codes.InvalidArgument, "new user role is required")
	}
	ur := domain.UserRole(req.GetNewUserRole())
	if ur == "" {
		return status.Error(codes.InvalidArgument, "new user role must be 'admin', 'super_user', 'ordinary'")
	}
	if req.GetReqUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "req user id is required")
	}
	return nil
}
