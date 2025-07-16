package cmd

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/tests/stores"
)

func TestDeleteSubscriberEmptyWebAppId(tst *testing.T) {
	// Given
	id := ""
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewDeleteSubscriberCmd(webAppsStore, subscribersStore)

	// When
	err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestDeleteSubscriberWebAppNotFound(tst *testing.T) {
	// Given
	id := uuid.New().String()
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	cmd := cmd.NewDeleteSubscriberCmd(webAppsStore, subscribersStore)

	// When
	err := cmd.Execute(id, user)

	// Then
	if err != sql.ErrNoRows {
		tst.Errorf("Ended with invalid error: %s", err)
	}
}

func TestSuccessDeleteSubscriber(tst *testing.T) {
	// Given
	user := "test@test.pl"
	webAppsStore := stores.NewInMemWebAppsStore()
	subscribersStore := stores.NewInMemSubscribersStore()
	app := types.NewWebApp("Google", "https://google.com", types.Interval1H, user, 200)
	app, _ = webAppsStore.Insert(app, nil)
	subscriber := types.NewSubscriber(app.Id, user)
	subscribersStore.Insert(subscriber, nil)
	cmd := cmd.NewDeleteSubscriberCmd(webAppsStore, subscribersStore)

	// When
	err := cmd.Execute(app.Id, user)

	// Then
	if err != nil {
		tst.Errorf("Deleting subscriber failed: %s", err)
	}

	subscribers, err := subscribersStore.GetByWebAppIdOrderByEmailAsc(app.Id, nil)

	if err != nil {
		tst.Errorf("Fetching subscribers failed: %s", err)
	}

	if len(subscribers) != 0 {
		tst.Errorf("Subscriber not deleted")
	}
}
