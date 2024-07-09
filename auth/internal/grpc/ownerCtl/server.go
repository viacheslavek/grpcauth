package ownerCtl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models"
	"github.com/viacheslavek/grpcauth/auth/internal/domain/models/validator"
	"github.com/viacheslavek/grpcauth/auth/internal/lib/logger/sl"
	"github.com/viacheslavek/grpcauth/auth/internal/services/ownerCtl"
	"github.com/viacheslavek/grpcauth/auth/internal/storage"
)

type OwnerCtl interface {
	CreateOwner(ctx context.Context, owner models.Owner) error
	UpdateOwner(ctx context.Context, owner models.Owner) error
	DeleteOwner(ctx context.Context, owner models.Owner) error
	GetOwner(ctx context.Context, owner models.Owner) (models.Owner, error)

	LoginOwner(ctx context.Context, owner models.Owner, appId int) (token string, err error)
}

type serverAPI struct {
	authv1.UnimplementedOwnerControllerServer
	octl OwnerCtl
	lg   *slog.Logger
}

func Register(gRPC *grpc.Server, octl OwnerCtl, lg *slog.Logger) {
	authv1.RegisterOwnerControllerServer(gRPC, &serverAPI{octl: octl, lg: lg})
}

// CreateOwner Creates a user in the table by email, login, and password
func (s *serverAPI) CreateOwner(
	ctx context.Context, req *authv1.CreateOwnerRequest,
) (*authv1.Response, error) {
	const op = "auth.CreateOwner"

	o := models.Owner{}
	if err := o.SetEmail(req.GetEmail()); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set email %v", op, err))
	}
	if err := o.SetLogin(req.GetLogin()); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set login %v", op, err))
	}
	if err := o.SetPassword(req.GetPassword()); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set password %v", op, err))
	}

	if err := s.octl.CreateOwner(ctx, o); err != nil {
		s.lg.With(
			slog.String("op", op),
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
	const op = "auth.UpdateOwner"

	o := models.Owner{}
	if err := o.SetId(req.GetId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set id %v", op, err))
	}

	if err := o.SetEmail(req.GetEmail()); err != nil && !errors.Is(err, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set email %v", op, err))
	}
	if err := o.SetLogin(req.GetLogin()); err != nil && !errors.Is(err, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set login %v", op, err))
	}
	if err := o.SetPassword(req.GetPassword()); err != nil && !errors.Is(err, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set password %v", op, err))
	}

	if err := s.octl.UpdateOwner(ctx, o); err != nil {
		s.lg.With(
			slog.String("op", op),
		).Error("failed to update owner", sl.Err(err))

		if errors.Is(err, ownerCtl.ErrInvalidCredentials) {
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
	const op = "auth.DeleteOwner"

	o := models.Owner{}
	errIdVal := o.SetId(req.GetId())
	errLoginVal := o.SetLogin(req.GetLogin())

	if errors.Is(errIdVal, validator.ErrEmptyParameter) && errors.Is(errLoginVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, "empty all delete parameters")
	}
	if errIdVal != nil && !errors.Is(errIdVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set id %v", op, errIdVal))
	}
	if errLoginVal != nil && !errors.Is(errLoginVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set login %v", op, errIdVal))
	}

	if err := s.octl.DeleteOwner(ctx, o); err != nil {
		s.lg.With(
			slog.String("op", op),
		).Error("failed to delete owner", sl.Err(err))

		if errors.Is(err, ownerCtl.ErrInvalidCredentials) {
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
	const op = "papa auth.GetOwner"

	o := models.Owner{}
	errIdVal := o.SetId(req.GetId())
	errLoginVal := o.SetLogin(req.GetLogin())

	if errors.Is(errIdVal, validator.ErrEmptyParameter) && errors.Is(errLoginVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, "empty all get parameters")
	}
	if errIdVal != nil && !errors.Is(errIdVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set id %v", op, errIdVal))
	}
	if errLoginVal != nil && !errors.Is(errLoginVal, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set login %v", op, errLoginVal))
	}

	owner, err := s.octl.GetOwner(ctx, o)
	if err != nil {
		s.lg.With(
			slog.String("op", op),
		).Error("failed to get owner", sl.Err(err))

		if errors.Is(err, ownerCtl.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or id")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.Owner{
		Id: owner.Id(), Email: owner.Email(), Login: owner.Login(), PasswordHash: string(owner.PassHash()),
	}, nil
}

// LoginOwner Issues a JWT token by login and password
func (s *serverAPI) LoginOwner(
	ctx context.Context, req *authv1.LoginOwnerRequest,
) (*authv1.LoginResponse, error) {
	const op = "auth.LoginOwner"
	if req.GetLogin() == "" || req.GetPassword() == "" || req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: empty parameters", op))
	}

	o := models.Owner{}
	if err := o.SetLogin(req.GetLogin()); err != nil && !errors.Is(err, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set login %v", op, err))
	}
	if err := o.SetPassword(req.GetPassword()); err != nil && !errors.Is(err, validator.ErrEmptyParameter) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%s: failed set password %v", op, err))
	}

	token, err := s.octl.LoginOwner(ctx, o, int(req.GetAppId()))
	if err != nil {
		s.lg.With(
			slog.String("op", op),
		).Error("failed to login owner", sl.Err(err))

		if errors.Is(err, ownerCtl.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{Token: token}, nil
}
