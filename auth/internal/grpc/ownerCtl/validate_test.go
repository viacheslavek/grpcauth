package ownerCtl

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email       string
		expectError bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"", true},
	}

	for _, test := range tests {
		err := validateEmail(test.email)
		if test.expectError && err == nil {
			t.Errorf("Expected error for email: %s, but got none", test.email)
		} else if !test.expectError && err != nil {
			t.Errorf("Did not expect error for email: %s, but got: %v", test.email, err)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password    string
		expectError bool
	}{
		{"password123", false},
		{"pass123", true},  // too short
		{"password", true}, // no digit
		{"12345678", true}, // no letter
		{"", true},         // empty
	}

	for _, test := range tests {
		err := validatePassword(test.password)
		if test.expectError && err == nil {
			t.Errorf("Expected error for password: %s, but got none", test.password)
		} else if !test.expectError && err != nil {
			t.Errorf("Did not expect error for password: %s, but got: %v", test.password, err)
		}
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		login       string
		expectError bool
	}{
		{"username123", false},
		{"user_name", true}, // invalid character
		{"username!", true}, // invalid character
		{"", true},          // empty
	}

	for _, test := range tests {
		err := validateLogin(test.login)
		if test.expectError && err == nil {
			t.Errorf("Expected error for login: %s, but got none", test.login)
		} else if !test.expectError && err != nil {
			t.Errorf("Did not expect error for login: %s, but got: %v", test.login, err)
		}
	}
}
