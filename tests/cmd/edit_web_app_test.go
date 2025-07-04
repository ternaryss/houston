package cmd

import (
	"database/sql"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/mocks"
	"github.com/ternaryss/houston/tests/stores"
)

func TestEditWebAppNoId(tst *testing.T) {
	// Given
	id := ""
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	form, err := types.NewWebAppForm(nil)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestEditWebAppNotFound(tst *testing.T) {
	// Given
	id := uuid.New().String()
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	form, err := types.NewWebAppForm(nil)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestEditWebAppEmptyName(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "")
	data.Set("url", "https://google.com")
	req, err := mocks.MockRequest("PUT", data)

	if err != nil {
		tst.Errorf("Edit web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, app.Id, user)

	if err != nil {
		tst.Errorf("Editing web application failed: %s", err)
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

func TestEditWebAppEmptyUrl(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "")
	req, err := mocks.MockRequest("PUT", data)

	if err != nil {
		tst.Errorf("Edit web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, app.Id, user)

	if err != nil {
		tst.Errorf("Editing web application failed: %s", err)
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
func TestEditWebAppInvalidUrl(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	req, err := mocks.MockRequest("PUT", data)

	if err != nil {
		tst.Errorf("Edit web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, app.Id, user)

	if err != nil {
		tst.Errorf("Editing web application failed: %s", err)
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

func TestSuccessEditWebApp(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewEditWebAppCmd(webAppsStore)
	data := url.Values{}
	data.Set("name", "Yahoo")
	data.Set("url", "https://yahoo.com")
	req, err := mocks.MockRequest("PUT", data)

	if err != nil {
		tst.Errorf("Edit web app request creation failed: %s", err)
	}

	form, err := types.NewWebAppForm(req)

	if err != nil {
		tst.Errorf("Edit web app form creation failed: %s", err)
	}

	// When
	err = cmd.Execute(form, app.Id, user)

	if err != nil {
		tst.Errorf("Editing web application failed: %s", err)
	}

	// Then
	app, err = webAppsStore.GetByIdAndUserEmail(app.Id, user)

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
