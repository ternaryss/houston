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
	subscribersStore := stores.NewInMemSubscribersStore()
	healthChecksStore := stores.NewInMemHealthChecksStore()
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore, subscribersStore, healthChecksStore)

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
	subscribersStore := stores.NewInMemSubscribersStore()
	healthChecksStore := stores.NewInMemHealthChecksStore()
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore, subscribersStore, healthChecksStore)

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
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	healthChecksStore := stores.NewInMemHealthChecksStore()
	app, _ = webAppsStore.Insert(app)
	subscriber := types.NewSubscriber(app.Id, user)
	subscribersStore.Insert(subscriber)
	health := types.NewHealthCheck(app.Id, 200)
	healthChecksStore.Insert(health)
	cmd := cmd.NewDeleteWebAppCmd(webAppsStore, subscribersStore, healthChecksStore)

	// When
	err := cmd.Execute(app.Id, user)

	// Then
	if err != nil {
		tst.Errorf("Deleting web application failed: %s", err)
	}

	if _, err := webAppsStore.GetByIdAndUserEmail(app.Id, user); err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}

	subscribers, err := subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id)

	if err != nil {
		tst.Errorf("Fetching subscribers failed: %s", err)
	}

	if len(subscribers) != 0 {
		tst.Errorf("Subscribers not deleted")
	}

	checks, err := healthChecksStore.GetByWebAppId(app.Id)

	if err != nil {
		tst.Errorf("Fetching health checks failed: %s", err)
	}

	if len(checks) != 0 {
		tst.Errorf("Health checks not deleted")
	}
}
