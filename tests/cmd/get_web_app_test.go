package cmd

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/stores"
)

func TestGetWebAppNoId(tst *testing.T) {
	// Given
	id := ""
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewGetWebAppCmd(webAppsStore, subscribersStore)

	// When
	_, _, err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestGetWebAppNotFound(tst *testing.T) {
	// Given
	id := uuid.New().String()
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewGetWebAppCmd(webAppsStore, subscribersStore)

	// When
	_, _, err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestSuccessGetWebApp(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	app, _ = webAppsStore.Insert(app, nil)
	subscribersStore := stores.NewInMemSubscribersStore()
	subscriber := types.NewSubscriber(app.Id, user)
	subscribersStore.Insert(subscriber, nil)
	cmd := cmd.NewGetWebAppCmd(webAppsStore, subscribersStore)

	// When
	savedApp, savedSubscribers, err := cmd.Execute(app.Id, user)

	// Then
	if err != nil {
		tst.Errorf("Fetching web application failed: %s", err)
	}

	if savedApp.Id == "" {
		tst.Error("Id is empty")
	}

	if savedApp.Name != app.Name {
		tst.Error("Name not matched")
	}

	if savedApp.Url != app.Url {
		tst.Error("URL not matched")
	}

	if savedApp.Status != app.Status {
		tst.Error("Status not matched")
	}

	if savedApp.Interval != app.Interval {
		tst.Error("Interval not matched")
	}

	if savedApp.UserEmail != app.UserEmail {
		tst.Error("User e-mail not matched")
	}

	if len(savedSubscribers) != 1 {
		tst.Error("Invalid subscribers size")
	}

	if savedSubscribers[0].WebAppId != savedApp.Id {
		tst.Error("Subscriber web app id not matched")
	}

	if savedSubscribers[0].Email != user {
		tst.Error("Subscriber e-mail not matched")
	}
}
