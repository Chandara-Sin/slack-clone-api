package store

import (
	"database/sql"
	"fmt"
	"slack-clone-api/config"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func CreateDB() *bun.DB {
	config.InitConfig()
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("app.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.dbname"),
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
