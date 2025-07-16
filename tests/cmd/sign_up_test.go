package cmd

import (
	"net/url"
	"testing"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/mocks"
	"github.com/ternaryss/houston/tests/stores"
)

func TestSignUpEmptyEmail(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "")
	data.Set("password", "1qaz@WSX3edc")
	data.Set("repeatPassword", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	emailErr, exists := form.Errors["email"]

	if !exists {
		tst.Error("No address e-mail validation error")
	}

	if emailErr.Field != "email" {
		tst.Error("Invalid field indication")
	}

	if emailErr.Msg != "Address e-mail is required." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpInvalidEmail(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test")
	data.Set("password", "1qaz@WSX3edc")
	data.Set("repeatPassword", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	emailErr, exists := form.Errors["email"]

	if !exists {
		tst.Error("No address e-mail validation error")
	}

	if emailErr.Field != "email" {
		tst.Error("Invalid field indication")
	}

	if emailErr.Msg != "Invalid address e-mail." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpEmptyPassword(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "")
	data.Set("repeatPassword", "")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	passwordErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No password validation error")
	}

	if passwordErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if passwordErr.Msg != "Password is required." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpPasswordLength(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "1234567")
	data.Set("repeatPassword", "")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	passwordErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No password validation error")
	}

	if passwordErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if passwordErr.Msg != "Invalid password (8 characters, one capital, one digit, one special character)." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpInvalidPassword(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "qazwsxedc")
	data.Set("repeatPassword", "")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	passwordErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No password validation error")
	}

	if passwordErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if passwordErr.Msg != "Invalid password (8 characters, one capital, one digit, one special character)." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpPasswordsNotMatch(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "1qaz@WSX3edc")
	data.Set("repeatPassword", "cde3XSW@zaq1")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	passwordErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No password validation error")
	}

	if passwordErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if passwordErr.Msg != "Passwords not match." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignUpUserExists(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	user := types.NewUser("test@test.com", "1qaz@WSX3edc")

	if _, err := usersStore.Insert(user, nil); err != nil {
		tst.Errorf("User can not be saved: %s", err)
	}

	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "1qaz@WSX3edc")
	data.Set("repeatPassword", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	signUpErr, exists := form.Errors["email"]

	if !exists {
		tst.Error("No sign up validation error")
	}

	if signUpErr.Field != "email" {
		tst.Error("Invalid field indication")
	}

	if signUpErr.Msg != "Already in use." {
		tst.Error("Invalid validation error message")
	}
}

func TestSuccessSignUp(tst *testing.T) {
	// Given
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignUpCmd(usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "1qaz@WSX3edc")
	data.Set("repeatPassword", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	form, err := types.NewSignUpForm(req)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	// Then
	user, err := usersStore.GetByEmail(form.Email, nil)

	if err != nil {
		tst.Errorf("Reading user data failed: %s", err)
	}

	if user.Email != form.Email {
		tst.Error("Address e-mail not matched")
	}

	if user.Password == "" {
		tst.Error("Password not set for user")
	}
}
