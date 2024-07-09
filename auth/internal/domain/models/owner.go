package models

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/viacheslavek/grpcauth/auth/internal/domain/models/validator"
)

type Owner struct {
	id       int64
	email    string
	login    string
	password string
	passHash []byte
}

type OwnerKey struct {
	Id    int64
	Login string
}

const emptyId = 0

func (o *Owner) SetId(id int64) error {
	if id == emptyId {
		return validator.ErrEmptyParameter
	}
	if id < 0 {
		return fmt.Errorf("id can't be less than zero, given %d", id)
	}

	o.id = id

	return nil
}

func (o *Owner) SetLogin(login string) error {
	if len(login) == 0 {
		return validator.ErrEmptyParameter
	}

	if err := validator.ValidateLogin(login); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	o.login = login

	return nil
}

func (o *Owner) SetEmail(email string) error {
	if len(email) == 0 {
		return validator.ErrEmptyParameter
	}

	if err := validator.ValidateEmail(email); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	o.email = email

	return nil
}

func (o *Owner) SetPassword(password string) error {
	if len(password) == 0 {
		return validator.ErrEmptyParameter
	}

	if err := validator.ValidatePassword(password); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	o.password = password

	return nil
}

func (o *Owner) SetPassHash(passHash []byte) {
	o.passHash = passHash
}

func (o *Owner) Id() int64 {
	return o.id
}

func (o *Owner) Login() string {
	return o.login
}

func (o *Owner) Email() string {
	return o.email
}

func (o *Owner) Password() string {
	return o.password
}

func (o *Owner) PassHash() []byte {
	return o.passHash
}
