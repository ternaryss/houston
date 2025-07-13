package cli

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/ternaryss/houston/internal/app/cron"
	"github.com/ternaryss/houston/internal/app/db"
	"github.com/ternaryss/houston/internal/app/handlers"
	icmd "github.com/ternaryss/houston/internal/app/handlers/cmd"
	"github.com/ternaryss/houston/internal/app/settings"
	"github.com/ternaryss/houston/internal/app/types"
	"github.com/ternaryss/houston/internal/app/web"
)

var rootCmd = &cobra.Command{
	Use:   "houston [command]",
	Short: "Houston we have (no) problem - periodic web apps health check",
	Long:  "Houston is a lightweight and efficient web application designed for monitoring the availability of other web applications. Built with simplicity in mind, it provides an easy-to-use alternative to complex monitoring solutions like Grafana.",
	Run: func(cmd *cobra.Command, ags []string) {
		settings := settings.LoadSettings()
		dbProvider := db.NewDbProvider(settings)
		defer dbProvider.CloseConnection()
		dbProvider.MigrateDatabase()
		usersStore := db.NewUsersStore(dbProvider.Db())
		webAppsStore := db.NewWebAppsStore(dbProvider.Db())
		subscribersStore := db.NewSubscribersStore(dbProvider.Db())
		healthChecksStore := db.NewHealthChecksStore(dbProvider.Db())
		errorsHandler := handlers.NewErrorsHandler()
		dashboardHandler := handlers.NewDashboardHandler(webAppsStore)
		webAppsHandler := handlers.NewWebAppsHandler(webAppsStore, subscribersStore)
		healthChecksHandler := handlers.NewHealthChecksHandler(settings, webAppsStore, healthChecksStore)
		usersHandler := handlers.NewUsersHandler(settings, usersStore)
		retentionScheduler := cron.NewRetentionScheduler(settings, healthChecksHandler)
		fiveMinutesInterval := cron.NewFiveMinutesIntervalScheduler(healthChecksHandler)
		fifteenMinutesInterval := cron.NewFifteenMinutesIntervalScheduler(healthChecksHandler)
		oneHourInterval := cron.NewOneHourIntervalScheduler(healthChecksHandler)
		retentionScheduler.Run()
		fiveMinutesInterval.Run()
		fifteenMinutesInterval.Run()
		oneHourInterval.Run()
		server := web.NewServer(settings, errorsHandler, dashboardHandler, webAppsHandler, usersHandler)
		server.Run()
	},
}

var createUserCmd = &cobra.Command{
	Use:   "createuser [email] [password]",
	Short: "Backoffice user creation",
	Long:  "Backoffice user creation that is especially usable in environment where 'Sign up' feature must be disabled.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, ags []string) error {
		settings := settings.LoadSettings()
		dbProvider := db.NewDbProvider(settings)
		defer dbProvider.CloseConnection()
		dbProvider.MigrateDatabase()
		usersStore := db.NewUsersStore(dbProvider.Db())
		signUpCmd := icmd.NewSignUpCmd(usersStore)
		payload := url.Values{}
		payload.Set("email", ags[0])
		payload.Set("password", ags[1])
		payload.Set("repeatPassword", ags[1])
		req, _ := http.NewRequest("POST", "", bytes.NewBufferString(payload.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		form, _ := types.NewSignUpForm(req)
		signUpCmd.Execute(form)

		if len(form.Errors) > 0 {
			return fmt.Errorf("user creation failed - %s", form.Errors)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)
}

func Execute() {
	rootCmd.Execute()
}
