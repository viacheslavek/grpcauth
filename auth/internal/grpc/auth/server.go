package auth

import (
	"context"
	"google.golang.org/grpc"

	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
)

type serverAPI struct {
	authv1.UnimplementedOwnerControllerServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterOwnerControllerServer(gRPC, &serverAPI{})
}

// CreateOwner TODO: реализовать
func (s *serverAPI) CreateOwner(
	ctx context.Context, req *authv1.CreateOwnerRequest,
) (*authv1.Response, error) {
	return &authv1.Response{Code: 0, Message: "norm"}, nil
}

// UpdateOwner TODO: реализовать
func (s *serverAPI) UpdateOwner(
	context.Context, *authv1.UpdateOwnerRequest,
) (*authv1.Response, error) {
	panic("implement me")
}

// DeleteOwner TODO: реализовать
func (s *serverAPI) DeleteOwner(
	context.Context, *authv1.DeleteOwnerRequest,
) (*authv1.Response, error) {
	panic("implement me")
}

// GetOwner TODO: реализовать
func (s *serverAPI) GetOwner(
	context.Context, *authv1.GetOwnerRequest,
) (*authv1.Owner, error) {
	panic("implement me")
}

// LoginOwner TODO: реализовать
func (s *serverAPI) LoginOwner(
	context.Context, *authv1.LoginOwnerRequest,
) (*authv1.LoginResponse, error) {
	panic("implement me")
}

// TODO: комментарии к методам добавить
