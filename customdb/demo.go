package custom_db_demo

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go-sql-log/logger"
	"go-sql-log/types"
)

func Run() {
	dsn := ":memory:"
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db")
	}

	//db.MustExec("CREATE table t (i integer, text varchar(16))")
	// 底层调用 MustExec(db, query, args...)
	// db 实现 sqlx.Execer 接口 和 sqlx.Queryer 接口 等

	newDB := &QueryLogger{
		db:     db,
		logger: log.Logger,
	}

	// 后续需要使用 newDB 替代 db，使用 sqlx 代替 sql
	sqlx.MustExec(newDB, "CREATE table t (id integer primary key , text varchar(16))")
	sqlx.MustExec(newDB, "insert into t (text) values(?),(?)", "foo", "bar")

	// 查询
	list := make([]*types.Table, 0)
	err = sqlx.Select(newDB, &list, "SELECT * FROM t")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to query")
	}
	log.Info().Any("list", list).Msg("select")
}

type SQLXDBInterface interface {
	sqlx.Execer
	sqlx.Queryer
	// ...
}

type QueryLogger struct {
	db     SQLXDBInterface
	logger zerolog.Logger
}

var _ sqlx.Execer = (SQLXDBInterface)(nil)

func (q *QueryLogger) Exec(query string, args ...interface{}) (sql.Result, error) {
	q.logger.Info().Str("sql", logger.ExplainSQL(
		query,
		nil,
		`'`,
		args...,
	)).Msg("exec")
	return q.db.Exec(query, args...)
}

func (q *QueryLogger) Query(query string, args ...interface{}) (*sql.Rows, error) {
	q.logger.Info().Str("sql", logger.ExplainSQL(
		query,
		nil,
		`'`,
		args...,
	)).Msg("query")
	return q.db.Query(query, args...)
}

func (q *QueryLogger) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	q.logger.Info().Str("sql", logger.ExplainSQL(
		query,
		nil,
		`'`,
		args...,
	)).Msg("queryx")
	return q.db.Queryx(query, args...)
}

func (q *QueryLogger) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	q.logger.Info().Str("sql", logger.ExplainSQL(
		query,
		nil,
		`'`,
		args...,
	)).Msg("queryrowx")
	return q.db.QueryRowx(query, args...)
}
