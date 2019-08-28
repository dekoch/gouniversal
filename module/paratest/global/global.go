package global

import "github.com/dekoch/gouniversal/shared/io/sqlite3"

var (
	DBConn        sqlite3.SQLite
	LytAuftrag    sqlite3.Layout
	LytParameter  sqlite3.Layout
	LytEinzelPara sqlite3.Layout
)
