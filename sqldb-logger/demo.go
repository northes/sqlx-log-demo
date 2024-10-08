package sqldb_logger_demo

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	dsn := ":memory:"
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to ping db")
	}

	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqlx.NewDb(sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter), db.DriverName())

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to ping db")
	}

	_ = db.MustExec("CREATE table t (id integer primary key , text varchar(16))")
	_ = db.MustExec("insert into t (text) values(?),(?)", "foo", "bar")
	_, err = db.Query("select * from t")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to query")
	}
}
