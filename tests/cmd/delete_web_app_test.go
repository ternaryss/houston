package cmd

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/stores"
)

func TestDeleteWebAppNoId(tst *testing.T) {
	// Given
	id := ""
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore)

	// When
	err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestDeleteWebAppNotFound(tst *testing.T) {
	// Given
	id := uuid.New().String()
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore)

	// When
	err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestSuccessDeleteWebApp(tst *testing.T) {
	// Given
	user := "test@test.pl"
	app := types.NewWebApp("Google", "https://google.com", user)
	webAppsStore := stores.NewInMemWebAppsStore()
	webAppsStore.Insert(app)
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore)

	// When
	err := cmd.Execute(app.Id, user)

	// Then
	if err != nil {
		tst.Errorf("Deleting web application failed: %s", err)
	}

	if _, err := webAppsStore.GetByIdAndUserEmail(app.Id, user); err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}
