package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	id := int64(999999)
	newLogin := gofakeit.Username()
	newEmail, errGVE := generateValidEmail(1000)
	assert.NoError(t, errGVE, "email generate failed")

	newPassword := generateValidPassword()

	_, err := s.OwnerClient.UpdateOwner(s.Ctx, &authv1.UpdateOwnerRequest{
		Id:       id,
		Login:    newLogin,
		Email:    newEmail,
		Password: newPassword,
	})
	require.Error(t, err, "expected error when updating non-existent owner")
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code(), "expected status code InvalidArgument")
	assert.Contains(t, st.Message(), "invalid id", "expected owner not found message")
}

func TestDeleteOwner_Failures(t *testing.T) {
	s := suite.New(t)

	// Try deleting a non-existent owner
	id := int64(99999) // Assume this ID does not exist

	_, err := s.OwnerClient.DeleteOwner(s.Ctx, &authv1.DeleteOwnerRequest{
		Id: id,
	})

	require.Error(t, err, "expected error when deleting non-existent owner")
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code(), "expected status code InvalidArgument")
	assert.Contains(t, st.Message(), "invalid login or id", "expected owner not found message")
}

func TestGetOwner_Failures(t *testing.T) {
	s := suite.New(t)

	// Try getting a non-existent owner
	login := "nonExistentLoginForTest"

	_, err := s.OwnerClient.GetOwner(s.Ctx, &authv1.GetOwnerRequest{
		Login: login,
	})

	require.Error(t, err, "expected error when getting non-existent owner")
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code(), "expected status code InvalidArgument")
	assert.Contains(t, st.Message(), "invalid login or id", "expected owner not found message")
}
