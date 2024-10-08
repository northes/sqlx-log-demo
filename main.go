package main

import (
	customdbdemo "go-sql-log/customdb"
	"go-sql-log/logger"
	sqldbloggerdemo "go-sql-log/sqldb-logger"
	sqlhooksdemo "go-sql-log/sqlhooks"
)

func init() {
	logger.Init()
}

func main() {
	// 在sqlx的层面实现，通过实现 sqlx 的相关接口，来达到打印日志的目的，相对来说更复杂（层级更高）
	customdbdemo.Run()

	// 在driver层面实现，需要传入logger
	sqldbloggerdemo.Run()

	// 在driver的层面实现（层级更低），实现Before和After即可，相对更加灵活
	sqlhooksdemo.Run()
}
