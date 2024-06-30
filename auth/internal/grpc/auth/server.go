package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
)

type Auth interface {
	CreateOwner(ctx context.Context, owner models.Owner) error
	UpdateOwner(ctx context.Context, owner models.Owner) error
	DeleteOwner(ctx context.Context, owner models.Owner) error
	GetOwner(ctx context.Context, owner models.Owner) (models.Owner, error)

	LoginOwner(ctx context.Context, owner models.Owner, appId int) (token string, err error)
}

type serverAPI struct {
	authv1.UnimplementedOwnerControllerServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterOwnerControllerServer(gRPC, &serverAPI{auth: auth})
}

const emptyId = 0

// CreateOwner TODO comment
func (s *serverAPI) CreateOwner(
	ctx context.Context, req *authv1.CreateOwnerRequest,
) (*authv1.Response, error) {

	if err := validateEmail(req.GetEmail()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validateLogin(req.GetLogin()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validatePassword(req.GetPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: правлю models.NewOwner()
	if err := s.auth.CreateOwner(ctx, models.NewOwner()); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success create owner"}, nil
}

// UpdateOwner TODO comment
func (s *serverAPI) UpdateOwner(
	ctx context.Context, req *authv1.UpdateOwnerRequest,
) (*authv1.Response, error) {

	if req.GetId() == emptyId {
		return nil, status.Error(codes.InvalidArgument, "empty id")
	}

	if req.GetEmail() == "" && req.GetLogin() == "" && req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "all update params is empty")
	}

	if req.GetEmail() != "" && validateEmail(req.GetEmail()) != nil {
		return nil, status.Error(codes.InvalidArgument, validateEmail(req.GetEmail()).Error())
	}

	if req.GetLogin() != "" && validateLogin(req.GetLogin()) != nil {
		return nil, status.Error(codes.InvalidArgument, validateLogin(req.GetLogin()).Error())
	}

	if req.GetPassword() != "" && validatePassword(req.GetPassword()) != nil {
		return nil, status.Error(codes.InvalidArgument, validatePassword(req.GetPassword()).Error())
	}

	// TODO: правлю models.NewOwner()
	if err := s.auth.UpdateOwner(ctx, models.NewOwner()); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success update owner"}, nil
}

// DeleteOwner TODO comment
func (s *serverAPI) DeleteOwner(
	ctx context.Context, req *authv1.DeleteOwnerRequest,
) (*authv1.Response, error) {
	if req.GetId() == emptyId && req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delete parameters")
	}

	// TODO: правлю models.NewOwner()
	if err := s.auth.DeleteOwner(ctx, models.NewOwner()); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success delete owner"}, nil
}

// GetOwner TODO comment
func (s *serverAPI) GetOwner(
	ctx context.Context, req *authv1.GetOwnerRequest,
) (*authv1.Owner, error) {
	if req.GetId() == emptyId && req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty get parameters")
	}

	// TODO: правлю models.NewOwner()
	owner, err := s.auth.GetOwner(ctx, models.NewOwner())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	// TODO: заполняю owner через owner
	fmt.Println(owner)
	return &authv1.Owner{}, nil
}

// LoginOwner TODO comment
func (s *serverAPI) LoginOwner(
	ctx context.Context, req *authv1.LoginOwnerRequest,
) (*authv1.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" || req.GetAppId() == emptyId {
		return nil, status.Error(codes.InvalidArgument, "empty login parameters")
	}

	// TODO: правлю models.NewOwner()
	token, err := s.auth.LoginOwner(ctx, models.NewOwner(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Code: int32(codes.OK), Token: token}, nil
}
