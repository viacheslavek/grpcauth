package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"

	"github.com/viacheslavek/grpcauth/auth/tests/suite"
)

func TestCreateOwner_Failures(t *testing.T) {
	s := suite.New(t)

	login := gofakeit.Username()
	email, errGVE := generateValidEmail(1000)
	assert.NoError(t, errGVE, "email generate failed")
	password := generateValidPassword()

	createOwnerAndCheckSuccess(s, t, login, email, password)

	// Try to create existent owner
	_, err := s.OwnerClient.CreateOwner(s.Ctx, &authv1.CreateOwnerRequest{
		Login:    login,
		Email:    email,
		Password: password,
	})

	require.Error(t, err, "expected error when creating owner with existing email or login")

	st, ok := status.FromError(err)
	assert.True(t, ok, "expected gRPC status error")
	assert.Equal(t, codes.AlreadyExists, st.Code(), "expected AlreadyExists error code")

	assert.Contains(t, st.Message(),
		"user already exists", "expected error message to contain 'user already exists'")
}

func TestUpdateOwner_Failures(t *testing.T) {
	s := suite.New(t)

	// Try updating a non-existent owner
	id := int64(99999) // Assume this ID does not exist
	newLogin := "new_login"
	newEmail := "new_email@example.com"
	newPassword := generateValidPassword()

	_, err := s.OwnerClient.UpdateOwner(s.Ctx, &authv1.UpdateOwnerRequest{
		Id:       id,
		Login:    newLogin,
		Email:    newEmail,
		Password: newPassword,
	})
	require.Error(t, err, "expected error when updating non-existent owner")
	assert.Contains(t, err.Error(), "owner not found", "expected owner not found error")
}

func TestDeleteOwner_Failures(t *testing.T) {
	s := suite.New(t)

	// Try deleting a non-existent owner
	id := int64(99999) // Assume this ID does not exist

	_, err := s.OwnerClient.DeleteOwner(s.Ctx, &authv1.DeleteOwnerRequest{
		Id: id,
	})

	require.Error(t, err, "expected error when deleting non-existent owner")
	assert.Contains(t, err.Error(), "owner not found", "expected owner not found error")
}

func TestGetOwner_Failures(t *testing.T) {
	s := suite.New(t)

	// Try getting a non-existent owner
	login := "non_existent_login_for_test"

	_, err := s.OwnerClient.GetOwner(s.Ctx, &authv1.GetOwnerRequest{
		Login: login,
	})

	require.Error(t, err, "expected error when getting non-existent owner")
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code(), "expected status code InvalidArgument")
	assert.Contains(t, st.Message(), "invalid login or id", "expected owner not found message")
}
