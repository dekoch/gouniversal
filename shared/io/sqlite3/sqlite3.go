package sqlite3

import (
	"database/sql"
	"errors"
	"strings"
	"sync"

	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/sbool"

	_ "github.com/mattn/go-sqlite3"
)

type DataType string

const (
	TypeTEXT    DataType = "TEXT"
	TypeNUMERIC DataType = "NUMERIC"
	TypeINTEGER DataType = "INTEGER"
	TypeREAL    DataType = "REAL"
	TypeBLOB    DataType = "BLOB"
)

type SQLite struct {
	DB       *sql.DB
	Tx       *sql.Tx
	dbIsOpen sbool.Sbool
}

var mut sync.Mutex

func (sq *SQLite) Open(dataSourceName string) error {

	if functions.IsEmpty(dataSourceName) {
		return errors.New("empty dataSourceName")
	}

	mut.Lock()
	defer mut.Unlock()

	if sq.dbIsOpen.IsSet() {
		return errors.New("DB already open")
	}

	err := functions.CreateDir(dataSourceName)
	if err != nil {
		return err
	}

	sq.DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	err = sq.setOptions()
	if err != nil {
		return err
	}

	sq.dbIsOpen.Set()

	return nil
}

func (sq *SQLite) Close() error {

	mut.Lock()
	defer mut.Unlock()

	if sq.dbIsOpen.IsSet() == false {
		return errors.New("DB already closed")
	}

	err := sq.DB.Close()
	if err != nil {
		return err
	}

	sq.dbIsOpen.UnSet()

	return nil
}

func (sq *SQLite) TableExists(tableName string) (bool, error) {

	if functions.IsEmpty(tableName) {
		return false, errors.New("empty tableName")
	}

	var (
		err  error
		rows *sql.Rows
		ret  bool
	)

	func() {

		for i := 0; i <= 2; i++ {

			switch i {
			case 0:
				rows, err = sq.DB.Query("PRAGMA table_info(" + tableName + ")")

			case 1:
				defer rows.Close()

			case 2:
				for rows.Next() {

					ret = true
					return
				}
			}

			if err != nil {
				return
			}
		}
	}()

	return ret, err
}

func (sq *SQLite) CreateTable(tableName, table string) error {

	if functions.IsEmpty(tableName) {
		return errors.New("empty tableName")
	}

	if functions.IsEmpty(table) {
		return errors.New("empty table")
	}

	_, err := sq.DB.Exec("CREATE TABLE IF NOT EXISTS `" + tableName + "` (" + table + ")")
	return err
}

func (sq *SQLite) CreateTableFromLayout(la Layout) error {

	if functions.IsEmpty(la.tableName) {
		return errors.New("empty tableName")
	}

	if len(la.fields) == 0 {
		return errors.New("no field added")
	}

	var table string

	// add fields
	for _, f := range la.fields {

		if functions.IsEmpty(f.name) {
			continue
		}

		table += "`" + f.name + "` " + string(f.dType)

		if f.pk {
			table += " PRIMARY KEY"
		}

		if f.ai {
			table += " AUTOINCREMENT"
		}

		table += ", "
	}

	table = strings.Trim(table, " ")
	table = strings.Trim(table, ",")

	return sq.CreateTable(la.tableName, table)
}

func (sq *SQLite) DropTable(tableName string) error {

	if functions.IsEmpty(tableName) {
		return errors.New("empty tableName")
	}

	_, err := sq.DB.Exec("DROP TABLE IF EXISTS `" + tableName + "`")
	return err
}

func (sq *SQLite) Lock() {

	mut.Lock()
}

func (sq *SQLite) Unlock() {

	mut.Unlock()
}

func (sq *SQLite) setOptions() error {

	var err error

	ddl := `
	PRAGMA temp_store = MEMORY;
	`

	func() {

		for i := 0; i <= 1; i++ {

			switch i {
			case 0:
				_, err = sq.DB.Exec(ddl)

			case 1:
				sq.DB.SetMaxOpenConns(1)
			}

			if err != nil {
				return
			}
		}
	}()

	return err
}

type field struct {
	name  string
	dType DataType
	pk    bool
	ai    bool
}

type Layout struct {
	tableName string
	fields    []field
}

func (la *Layout) SetTableName(name string) {
	la.tableName = name
}

func (la *Layout) AddField(name string, dtype DataType, primaryKey, autoIncr bool) {

	var f field
	f.name = name
	f.dType = dtype
	f.pk = primaryKey
	f.ai = autoIncr

	la.fields = append(la.fields, f)
}
