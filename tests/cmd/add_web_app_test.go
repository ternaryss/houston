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
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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

func TestAddWebAppStatusTooLow(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "0")
	data.Set("interval", types.Interval1H)
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
	statusErr, exists := form.Errors["status"]

	if !exists {
		tst.Error("No status validation error")
	}

	if statusErr.Field != "status" {
		tst.Error("Invalid field indication")
	}

	if statusErr.Msg != "Invalid HTTP status" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppStatusTooHigh(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "1000")
	data.Set("interval", types.Interval1H)
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
	statusErr, exists := form.Errors["status"]

	if !exists {
		tst.Error("No status validation error")
	}

	if statusErr.Field != "status" {
		tst.Error("Invalid field indication")
	}

	if statusErr.Msg != "Invalid HTTP status" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppEmptyInterval(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", "")
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
	intervalErr, exists := form.Errors["interval"]

	if !exists {
		tst.Error("No interval validation error")
	}

	if intervalErr.Field != "interval" {
		tst.Error("Invalid field indication")
	}

	if intervalErr.Msg != "Interval is required" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppInvalidInterval(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", "1M")
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
	intervalErr, exists := form.Errors["interval"]

	if !exists {
		tst.Error("No interval validation error")
	}

	if intervalErr.Field != "interval" {
		tst.Error("Invalid field indication")
	}

	if intervalErr.Msg != "Invalid interval" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppEmptyEmail(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", "")
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
	emailErr, exists := form.Errors["email0"]

	if !exists {
		tst.Error("No email validation error")
	}

	if emailErr.Field != "email0" {
		tst.Error("Invalid field indication")
	}

	if emailErr.Msg != "Address e-mail is required" {
		tst.Error("Invalid validation error message")
	}
}

func TestAddWebAppInvalidEmail(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", "test")
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
	emailErr, exists := form.Errors["email0"]

	if !exists {
		tst.Error("No email validation error")
	}

	if emailErr.Field != "email0" {
		tst.Error("Invalid field indication")
	}

	if emailErr.Msg != "Invalid address e-mail" {
		tst.Error("Invalid validation error message")
	}
}

func TestSuccessAddWebAppWithoutNotify(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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

	if app.Status != form.Status {
		tst.Error("Status not matched")
	}

	if app.Interval != form.Interval {
		tst.Error("Interval not matched")
	}

	if app.UserEmail != user {
		tst.Error("User not matched")
	}
}

func TestSuccessAddWebAppWithNotify(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewAddWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", "test@test.pl")
	data.Add("notify[]", "test@test.com")
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

	if app.Status != form.Status {
		tst.Error("Status not matched")
	}

	if app.Interval != form.Interval {
		tst.Error("Interval not matched")
	}

	if app.UserEmail != user {
		tst.Error("User not matched")
	}

	subscribers, err := subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id)

	if err != nil {
		tst.Errorf("Reading subscribers data failed: %s", err)
	}

	if subscribers[0].Email != "test@test.com" {
		tst.Error("First subscriber not matched")
	}

	if subscribers[1].Email != "test@test.pl" {
		tst.Error("Second subscriber not matched")
	}
}
