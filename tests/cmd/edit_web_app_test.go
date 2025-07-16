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
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
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
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
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
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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

func TestEditWebAppStatusTooLow(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "0")
	data.Set("interval", types.Interval1H)
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

func TestEditWebAppStatusTooHigh(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "1000")
	data.Set("interval", types.Interval1H)
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

func TestEditWebAppEmptyInterval(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", "")
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

func TestEditWebAppInvalidInterval(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", "1M")
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

func TestEditWebAppEmptyEmail(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", "")
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

func TestEditWebAppInvalidEmail(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Google")
	data.Set("url", "https://google.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", "test")
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

func TestSuccessEditWebAppWithoutNotify(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	webAppsStore.Insert(app, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Yahoo")
	data.Set("url", "https://yahoo.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
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
	app, err = webAppsStore.GetByIdAndUserEmail(app.Id, user, nil)

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

func TestSuccessEditWebAppWithNotify(tst *testing.T) {
	// Given
	user := "test@test.pl"
	email := "test@test.com"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	app, _ = webAppsStore.Insert(app, nil)
	subscriber := types.NewSubscriber(app.Id, user)
	subscribersStore.Insert(subscriber, nil)
	cmd := cmd.NewEditWebAppCmd(webAppsStore, subscribersStore)
	data := url.Values{}
	data.Set("name", "Yahoo")
	data.Set("url", "https://yahoo.com")
	data.Set("status", "200")
	data.Set("interval", types.Interval1H)
	data.Add("notify[]", email)
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
	app, err = webAppsStore.GetByIdAndUserEmail(app.Id, user, nil)

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

	subscribers, err := subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id, nil)

	if err != nil {
		tst.Errorf("Reading subscribers data failed: %s", err)
	}

	if len(subscribers) != 1 {
		tst.Error("Subscribers size not matched")
	}

	if subscribers[0].Email != email {
		tst.Error("Subscriber not matched")
	}
}
