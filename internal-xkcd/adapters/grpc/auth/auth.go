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

type Auth struct {
	AuthClient pb.AuthClient
	Ctx        context.Context
}

func NewAuth(cfgAuthGRPC *config.AuthGRPC, ctx context.Context) (*Auth, error) {

	conn, cErr := grpc.NewClient(fmt.Sprintf("%s:%s", cfgAuthGRPC.Host, cfgAuthGRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if cErr != nil {
		return nil, cErr
	}
	client := pb.NewAuthClient(conn)

	return &Auth{
		AuthClient: client,
		Ctx:        ctx,
	}, nil
}

func (a *Auth) Login(login string, password string) (string, error) {
	loginRes, lErr := a.AuthClient.Login(a.Ctx, &pb.LoginRequest{Login: login, Password: password})

	if lErr != nil {
		st, ok := status.FromError(lErr)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return "", lErr
			case codes.NotFound:
				return "", domain.ErrUserNotFound
			default:
				return "", fmt.Errorf("unknown error: %v", st.Message())
			}
		}
		return "", lErr
	}
	return loginRes.GetToken(), nil
}

func (a *Auth) Register(login string, password string, role domain.UserRole) (int64, error) {
	regRes, rErr := a.AuthClient.Register(a.Ctx, &pb.RegisterRequest{Login: login, Password: password, Role: string(role)})
	if rErr != nil {
		st, ok := status.FromError(rErr)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return 0, rErr
			case codes.AlreadyExists:
				return 0, domain.ErrUserAlreadyExist
			default:
				return 0, fmt.Errorf("unknown error: %v", st.Message())
			}
		}
		return 0, rErr
	}
	return regRes.GetUserId(), nil
}

func (a *Auth) UserID(token string) (int64, error) {

	role, rErr := a.AuthClient.UserID(a.Ctx, &pb.UserIDRequest{Token: token})

	if rErr != nil {
		st, ok := status.FromError(rErr)
		if ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
				return 0, domain.ErrTokenExpired
			case codes.Unauthenticated:
				return 0, domain.ErrTokenNotValid
			case codes.InvalidArgument:
				return 0, rErr
			default:
				return 0, fmt.Errorf("unknown error: %v", st.Message())
			}
		}
		return 0, rErr
	}

	return role.GetUserId(), nil
}
