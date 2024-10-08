package sqlhooks_demo

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/rs/zerolog/log"

	"go-sql-log/logger"
	"go-sql-log/types"
)

const (
	registerName string = "_sqlite3_with_hooks"
)

func Run() {
	sql.Register(registerName, sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &Hooks{}))

	db, err := sqlx.Connect(registerName, ":memory:")
	//db, err := sql.Open("sqlite3WithHooks", ":memory:")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open db")
	}

	// 建表
	_ = db.MustExec("CREATE table t (id integer primary key , text varchar(16))")
	// 插入
	_ = db.MustExec("insert into t (text) values(?),(?)", "foo", "bar")
	// 查询
	rows, err := db.Queryx("select * from t")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to query")
	}
	list := make([]types.Table, 0)
	for rows.Next() {
		var t types.Table
		err = rows.StructScan(&t)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to scan")
		}
		list = append(list, t)
	}
	log.Info().Any("list", list).Msg("queryx")

	// 查询
	list2 := make([]*types.Table, 0)
	err = db.Select(&list2, "select * from t")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to query")
	}
	log.Info().Any("list", list2).Msg("select")
}

var _ sqlhooks.Hooks = (*Hooks)(nil)

type Hooks struct {
}

func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	//log.Info().Str("query", query).Any("args", args).Msg("before query")
	//log.Info().Msg("before query")
	return context.WithValue(ctx, "begin", time.Now()), nil
}

func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	//log.Info().Str("query", query).Any("args", args).Dur("elapsed", time.Since(begin)).Msg("after query")
	log.Debug().
		Str("query",
			logger.ExplainSQL(query, nil, `'`, args...),
		).
		Str("elapsed", time.Since(begin).String()).
		Msg("after query")

	// 尽量不要返回 error，此时 sql 语句已执行，但不会返回 result
	return ctx, nil
}
