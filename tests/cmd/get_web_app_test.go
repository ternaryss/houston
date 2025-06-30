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
	cmd := cmd.NewGetWebAppCmd(webAppsStore)

	// When
	_, err := cmd.Execute(id, user)

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
	cmd := cmd.NewGetWebAppCmd(webAppsStore)

	// When
	_, err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestSuccessGetWebApp(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewGetWebAppCmd(webAppsStore)

	// When
	result, err := cmd.Execute(app.Id, user)

	// Then
	if err != nil {
		tst.Errorf("Fetching web application failed: %s", err)
	}

	if result.Id == "" {
		tst.Error("Id is empty")
	}

	if result.Name != app.Name {
		tst.Error("Name not matched")
	}

	if result.Url != app.Url {
		tst.Error("URL not matched")
	}

	if result.UserEmail != app.UserEmail {
		tst.Error("User e-mail not matched")
	}
}
