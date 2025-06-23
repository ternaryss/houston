package cmd

import (
	"net/url"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/mocks"
	"github.com/ternaryss/houston/tests/stores"
)

func TestSignInEmptyEmail(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignInCmd(settings, usersStore)
	data := url.Values{}
	data.Set("email", "")
	data.Set("password", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(req)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	_, err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	fieldErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No validation error")
	}

	if fieldErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if fieldErr.Msg != "Invalid address e-mail or password." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignInEmptyPassword(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignInCmd(settings, usersStore)
	data := url.Values{}
	data.Set("email", "test@test.pl")
	data.Set("password", "")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(req)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	_, err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	fieldErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No validation error")
	}

	if fieldErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if fieldErr.Msg != "Invalid address e-mail or password." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignInInvalidEmail(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignInCmd(settings, usersStore)
	data := url.Values{}
	data.Set("email", "test")
	data.Set("password", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(req)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	_, err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	fieldErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No validation error")
	}

	if fieldErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if fieldErr.Msg != "Invalid address e-mail or password." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignInUserNotFound(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	cmd := cmd.NewSignInCmd(settings, usersStore)
	data := url.Values{}
	data.Set("email", "test@test.com")
	data.Set("password", "1qaz@WSX3edc")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(req)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	_, err = cmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	fieldErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No validation error")
	}

	if fieldErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if fieldErr.Msg != "Invalid address e-mail or password." {
		tst.Error("Invalid validation error message")
	}
}

func TestSignInInvalidPassword(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	signUpCmd := cmd.NewSignUpCmd(usersStore)
	signUpData := url.Values{}
	signUpData.Set("email", "test@test.com")
	signUpData.Set("password", "1qaz@WSX3edc")
	signUpData.Set("repeatPassword", "1qaz@WSX3edc")
	signUpReq, err := mocks.MockRequest("POST", signUpData)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	signUpForm, err := types.NewSignUpForm(signUpReq)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	signInCmd := cmd.NewSignInCmd(settings, usersStore)
	signInData := url.Values{}
	signInData.Set("email", "test@test.com")
	signInData.Set("password", "!QAZ2wsx#EDC")
	signInReq, err := mocks.MockRequest("POST", signInData)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(signInReq)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	err = signUpCmd.Execute(signUpForm)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	_, err = signInCmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	fieldErr, exists := form.Errors["password"]

	if !exists {
		tst.Error("No validation error")
	}

	if fieldErr.Field != "password" {
		tst.Error("Invalid field indication")
	}

	if fieldErr.Msg != "Invalid address e-mail or password." {
		tst.Error("Invalid validation error message")
	}
}

func TestSuccessSignIn(tst *testing.T) {
	// Given
	settings := settings.LoadSettings("../configs/app.yml")
	usersStore := stores.NewInMemUsersStore()
	signUpCmd := cmd.NewSignUpCmd(usersStore)
	user := "test@test.com"
	signUpData := url.Values{}
	signUpData.Set("email", user)
	signUpData.Set("password", "1qaz@WSX3edc")
	signUpData.Set("repeatPassword", "1qaz@WSX3edc")
	signUpReq, err := mocks.MockRequest("POST", signUpData)

	if err != nil {
		tst.Errorf("Sign up request creation failed: %s", err)
	}

	signUpForm, err := types.NewSignUpForm(signUpReq)

	if err != nil {
		tst.Errorf("Sign up form creation failed: %s", err)
	}

	signInCmd := cmd.NewSignInCmd(settings, usersStore)
	signInData := url.Values{}
	signInData.Set("email", user)
	signInData.Set("password", "1qaz@WSX3edc")
	signInReq, err := mocks.MockRequest("POST", signInData)

	if err != nil {
		tst.Errorf("Sign in request creation failed: %s", err)
	}

	form, err := types.NewSignInForm(signInReq)

	if err != nil {
		tst.Errorf("Sign in form creation failed: %s", err)
	}

	// When
	err = signUpCmd.Execute(signUpForm)

	if err != nil {
		tst.Errorf("Sign up command execution failed: %s", err)
	}

	result, err := signInCmd.Execute(form)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	// Then
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		result,
		claims,
		func(tkn *jwt.Token) (any, error) {
			return []byte(settings.Authorization.Secret), nil
		},
	)

	if err != nil {
		tst.Errorf("Sign in command execution failed: %s", err)
	}

	if !token.Valid {
		tst.Error("Access token is not valid")
	}

	if user != claims.Subject {
		tst.Error("Invalid user in access token")
	}
}
