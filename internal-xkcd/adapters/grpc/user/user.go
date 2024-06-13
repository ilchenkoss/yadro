package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	pb "myapp/pkg/proto/gen"
)

type User struct {
	UserClient pb.UserClient
	Ctx        context.Context
}

func NewUser(cfgAuthGRPC *config.AuthGRPC, ctx context.Context) (*User, error) {

	conn, cErr := grpc.NewClient(fmt.Sprintf("%s:%s", cfgAuthGRPC.Host, cfgAuthGRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if cErr != nil {
		return nil, cErr
	}
	client := pb.NewUserClient(conn)

	return &User{
		UserClient: client,
		Ctx:        ctx,
	}, nil
}

func (a *User) UserRole(userID int64) (domain.UserRole, error) {
	role, rErr := a.UserClient.UserRole(a.Ctx, &pb.UserRoleRequest{UserId: userID})

	if rErr != nil {
		st, ok := status.FromError(rErr)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return "", domain.ErrUserNotFound
			case codes.InvalidArgument:
				return "", rErr
			default:
				return "", fmt.Errorf("unknown error: %v", st.Message())
			}
		}
		return "", rErr
	}

	userRole := role.GetUserRole()

	ur := domain.UserRole(userRole)

	if ur != "" {
		return ur, nil
	}

	return "", domain.ErrUserRoleUnexpected
}

func (a *User) ChangeRole(reqUserID int64, newRole domain.UserRole) error {

	_, rErr := a.UserClient.ChangeRole(a.Ctx,
		&pb.ChangeRoleRequest{ReqUserId: reqUserID,
			NewUserRole: string(newRole)})

	if rErr != nil {
		st, ok := status.FromError(rErr)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return domain.ErrUserNotFound
			default:
				return fmt.Errorf("unknown error: %v", st.Message())
			}
		}
		return rErr

	}

	return nil
}
