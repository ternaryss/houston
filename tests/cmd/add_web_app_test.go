package cmd

import (
	"net/url"
	"testing"

	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/mocks"
	"github.com/ternaryss/houston/tests/stores"
)

func TestAddWebAppEmptyName(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "")
	data.Set("url", "https://google.com")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Add web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Add web app form creation failed: %s", err)
	}

	// When
	if _, err := cmd.Execute(form, user); err != nil {
		tst.Errorf("Add web app command execution failed: %s", err)
	}

	// Then
	nameErr, exists := form.Errors["name"]

	if !exists {
		tst.Error("No name validation error")
	}

	if nameErr.Field != "name" {
		tst.Error("Invalid field indication")
	}

	if nameErr.Msg != "Name is required" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppEmptyUrl(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Add web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Add web app form creation failed: %s", err)
	}

	// When
	if _, err := cmd.Execute(form, user); err != nil {
		tst.Errorf("Add web app command execution failed: %s", err)
	}

	// Then
	urlErr, exists := form.Errors["url"]

	if !exists {
		tst.Error("No url validation error")
	}

	if urlErr.Field != "url" {
		tst.Error("Invalid field indication")
	}

	if urlErr.Msg != "URL is required" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppInvalidUrl(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Add web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Add web app form creation failed: %s", err)
	}

	// When
	if _, err := cmd.Execute(form, user); err != nil {
		tst.Errorf("Add web app command execution failed: %s", err)
	}

	// Then
	urlErr, exists := form.Errors["url"]

	if !exists {
		tst.Error("No url validation error")
	}

	if urlErr.Field != "url" {
		tst.Error("Invalid field indication")
	}

	if urlErr.Msg != "Invalid URL" {
		tst.Error("Invalid validation error message")
	}
}

func TestSuccessAddWebApp(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	req, err := mocks.MockRequest("POST", data)

	if err != nil {
		tst.Errorf("Add web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Add web app form creation failed: %s", err)
	}

	// When
	id, err := cmd.Execute(form, user)

	if err != nil {
		tst.Errorf("Add web app command execution failed: %s", err)
	}

	// Then
	app, err := webAppsStore.GetByIdAndUserEmail(id, user)

	if err != nil {
		tst.Errorf("Reading web app data failed: %s", err)
	}

	if app.Name != form.Name {
		tst.Error("Name not matched")
	}

	if app.Url != form.Url {
		tst.Error("URL not matched")
	}

	if app.UserEmail != user {
		tst.Error("User not matched")
	}
}
