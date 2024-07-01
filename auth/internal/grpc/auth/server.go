package auth

import (
	"context"
	"errors"
	"log/slog"

	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger/sl"
	"github.com/viacheslavek/grpcauth/auth/internal/services/auth"
	"github.com/viacheslavek/grpcauth/auth/internal/storage"
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

		if errors.Is(err, storage.ErrOwnerExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Message: "Success create owner"}, nil
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

		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Message: "Success update owner"}, nil
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

		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or id")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Response{Message: "Success delete owner"}, nil
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

		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or id")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Owner{
		Id: owner.Id, Email: owner.Email, Login: owner.Login, PasswordHash: string(owner.PassHash),
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

		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		if errors.Is(err, auth.ErrInvalidApp) {
			return nil, status.Error(codes.InvalidArgument, "invalid app id")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}
