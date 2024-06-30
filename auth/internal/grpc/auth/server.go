package auth

import (
	"context"
	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
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
	lg   *slog.Logger
}

func Register(gRPC *grpc.Server, auth Auth, lg *slog.Logger) {
	authv1.RegisterOwnerControllerServer(gRPC, &serverAPI{auth: auth, lg: lg})
}

const emptyId = 0

// CreateOwner Creates a user in the table by email, login, and password
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

	if err := s.auth.CreateOwner(
		ctx,
		*models.NewOwner(
			models.WithEmail(req.GetEmail()),
			models.WithLogin(req.GetLogin()),
			models.WithPassword(req.GetPassword()),
		),
	); err != nil {
		s.lg.With(
			slog.String("op", "auth.CreateOwner"),
		).Error("failed to create owner", sl.Err(err))

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success create owner"}, nil
}

// UpdateOwner Updates the user's login, email, or password in the table by ID
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

	if err := s.auth.UpdateOwner(
		ctx,
		*models.NewOwner(
			models.WithId(req.GetId()),
			models.WithEmail(req.GetEmail()),
			models.WithLogin(req.GetLogin()),
			models.WithPassword(req.GetPassword()),
		),
	); err != nil {
		s.lg.With(
			slog.String("op", "auth.UpdateOwner"),
		).Error("failed to update owner", sl.Err(err))

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success update owner"}, nil
}

// DeleteOwner Deletes a user from the table by ID or login
func (s *serverAPI) DeleteOwner(
	ctx context.Context, req *authv1.DeleteOwnerRequest,
) (*authv1.Response, error) {
	if req.GetId() == emptyId && req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delete parameters")
	}

	if err := s.auth.DeleteOwner(
		ctx,
		*models.NewOwner(
			models.WithId(req.GetId()),
			models.WithLogin(req.GetLogin()),
		),
	); err != nil {
		s.lg.With(
			slog.String("op", "auth.DeleteOwner"),
		).Error("failed to delete owner", sl.Err(err))

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Code: int32(codes.OK), Message: "Success delete owner"}, nil
}

// GetOwner Retrieves a user from the table by ID or login
func (s *serverAPI) GetOwner(
	ctx context.Context, req *authv1.GetOwnerRequest,
) (*authv1.Owner, error) {
	if req.GetId() == emptyId && req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty get parameters")
	}

	owner, err := s.auth.GetOwner(
		ctx,
		*models.NewOwner(
			models.WithId(req.GetId()),
			models.WithLogin(req.GetLogin()),
		),
	)
	if err != nil {
		s.lg.With(
			slog.String("op", "auth.GetOwner"),
		).Error("failed to get owner", sl.Err(err))

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Owner{
		Id: owner.Id, Email: owner.Email, Login: owner.Login, Password: string(owner.PassHash),
	}, nil
}

// LoginOwner Issues a JWT token by login and password
func (s *serverAPI) LoginOwner(
	ctx context.Context, req *authv1.LoginOwnerRequest,
) (*authv1.LoginResponse, error) {
	if req.GetLogin() == "" || req.GetPassword() == "" || req.GetAppId() == emptyId {
		return nil, status.Error(codes.InvalidArgument, "empty login parameters")
	}

	token, err := s.auth.LoginOwner(
		ctx,
		*models.NewOwner(
			models.WithLogin(req.GetLogin()),
			models.WithPassword(req.GetPassword()),
		),
		int(req.GetAppId()))
	if err != nil {
		s.lg.With(
			slog.String("op", "auth.LoginOwner"),
		).Error("failed to login owner", sl.Err(err))

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Code: int32(codes.OK), Token: token}, nil
}
