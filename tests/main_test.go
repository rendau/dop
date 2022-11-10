package tests

import (
	"context"
	"os"
	"testing"

	"github.com/rendau/dop/adapters/db/pg"
	"github.com/rendau/dop/adapters/logger/zap"
	"github.com/spf13/viper"
)

var (
	bgCtx = context.Background()
	app   = struct {
		lg *zap.St
		db *pg.St
	}{
		lg: zap.New("info", true),
	}
)

func TestMain(m *testing.M) {
	var err error

	viper.AutomaticEnv()

	viper.SetDefault("PG_DSN", "postgres://localhost/dop")

	app.db, err = pg.New(true, app.lg, pg.OptionsSt{
		Dsn: viper.GetString("PG_DSN"),
	})
	if err != nil {
		app.lg.Fatal(err)
	}

	os.Exit(m.Run())
}
