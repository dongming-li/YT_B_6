package database

import (
	"database/sql"
	"fmt"

	mysql "github.com/go-sql-driver/mysql"
)

const (
	driverName       = "mysql"
	localDatasource  = "dbu309ytb6:rA@12XF3@/db309ytb6?parseTime=true"
	remoteDatasource = "dbu309ytb6:rA@12XF3@tcp(mysql.cs.iastate.edu:3306)/db309ytb6?parseTime=true"
)

var db *sql.DB

// Error implements go's standard error and wraps mysql driver's error.
type Error struct {
	Method  string
	Message string
	Code    uint16
}

func (e Error) Error() string {
	return fmt.Sprintf("database(%s): [%d] %s", e.Method, e.Code, e.Message)
}

func newError(err error, method string) error {
	if dErr, ok := err.(*mysql.MySQLError); ok {
		return Error{
			Method:  method,
			Message: dErr.Message,
			Code:    dErr.Number,
		}
	}
	return Error{
		Method:  method,
		Message: err.Error(),
		Code:    0,
	}
}

// Open returns an sql.DB to use in a query
func Open(remote bool) error {
	datasourceName := localDatasource
	if remote {
		datasourceName = remoteDatasource
	}
	var err error

	db, err = sql.Open(driverName, datasourceName)
	if err != nil {
		return newError(err, "Open")
	}

	err = db.Ping()
	if err != nil {
		return newError(err, "Open")
	}

	return nil
}

// Query takes a query string and arguments and returns the query's result
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	if len(args) == 0 {
		rows, err = db.Query(query)
	} else {
		rows, err = db.Query(query, args...)
	}

	if err != nil {
		return nil, newError(err, "Query")
	}
	return rows, nil
}

// QueryRow expects to return one row from the given query
func QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	if len(args) == 0 {
		return db.QueryRow(query), nil
	}
	return db.QueryRow(query, args...), nil
}

// Exec takes a query string and arguments and returns the result
func Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	if len(args) == 0 {
		result, err = db.Exec(query)
	} else {
		result, err = db.Exec(query, args...)
	}

	if err != nil {
		return nil, newError(err, "Exec")
	}
	return result, nil
}
