package e2e

import (
	"os"
	"testing"

	"github.com/callmhejerry/booking-api/tests/helpers"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var App *helpers.TestApp

func TestMain(m *testing.M) {

	app, close := helpers.SetupApp()
	App = app
	defer close()

	code := m.Run()
	app.CleanUpTestDB()
	os.Exit(code)
}
