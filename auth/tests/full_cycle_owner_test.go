package tests

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authv1 "github.com/viacheslavek/grpcauth/api/gen/go/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/viacheslavek/grpcauth/auth/tests/suite"
)

const passwordDefaultLen = 10

func TestFullCycleOwner_HappyPath(t *testing.T) {
	s := suite.New(t)

	login := gofakeit.Username()
	email, errGVE := generateValidEmail(1000)
	assert.NoError(t, errGVE, "email generate failed")
	password := generateValidPassword()
	createOwnerAndCheckSuccess(s, t, login, email, password)

	ownerID := getOwnerAndCheckSuccess(s, t, login, email, password)

	fmt.Println("id", ownerID)

	newLogin := gofakeit.Username()
	newEmail, errGVEN := generateValidEmail(1000)
	assert.NoError(t, errGVEN, "email generate failed")
	newPassword := generateValidPassword()
	updateOwnerAndCheckSuccess(s, t, ownerID, newLogin, newEmail, newPassword)

	deleteOwnerAndCheckSuccess(s, t, ownerID)
}

func generateValidPassword() string {
	return gofakeit.Password(true, true, true, true, true, passwordDefaultLen) + "1"
}

func generateValidEmail(maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		email := gofakeit.Email()
		if err := validation.Validate(email, validation.Required, is.Email); err == nil {
			return email, nil
		}
	}
	return "", fmt.Errorf("failed to generate a valid email after %d attempts", maxRetries)
}

func createOwnerAndCheckSuccess(s *suite.Suite, t *testing.T, login, email, password string) {
	res, err := s.OwnerClient.CreateOwner(s.Ctx, &authv1.CreateOwnerRequest{
		Login:    login,
		Email:    email,
		Password: password,
	})

	require.NoError(t, err,
		fmt.Sprintf("failed with owner: mail:%s login:%s password:%s", email, login, password))
	assert.NotEmpty(t, res.GetMessage(), "create response message")
}

func getOwnerAndCheckSuccess(s *suite.Suite, t *testing.T, login, email, password string) int64 {
	res, errGO := s.OwnerClient.GetOwner(s.Ctx, &authv1.GetOwnerRequest{
		Login: login,
	})

	require.NoError(t, errGO, fmt.Sprintf("failed with owner: login %s", login))
	assert.Equal(t, login, res.GetLogin(), "owner login")
	assert.Equal(t, email, res.GetEmail(), "owner email")
	passwordHash, errGFP := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, errGFP, "get password hash")
	assert.NoError(t,
		bcrypt.CompareHashAndPassword(passwordHash, []byte(password)), "owner password hash")

	return res.Id
}

func updateOwnerAndCheckSuccess(
	s *suite.Suite, t *testing.T,
	id int64, newLogin, newEmail, newPassword string,
) {
	res, err := s.OwnerClient.UpdateOwner(s.Ctx, &authv1.UpdateOwnerRequest{
		Id:       id,
		Login:    newLogin,
		Email:    newEmail,
		Password: newPassword,
	})

	require.NoError(t, err, "failed update")
	assert.NotEmpty(t, res.GetMessage())

	getOwnerAndCheckSuccess(s, t, newLogin, newEmail, newPassword)
}

func deleteOwnerAndCheckSuccess(s *suite.Suite, t *testing.T, id int64) {
	res, err := s.OwnerClient.DeleteOwner(s.Ctx, &authv1.DeleteOwnerRequest{
		Id: id,
	})

	require.NoError(t, err, "failed delete")
	assert.NotEmpty(t, res.GetMessage())
}
