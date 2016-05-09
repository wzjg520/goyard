package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	dsn      string
	connection *sql.DB
}

type oneRow map[string]interface{}

func (db *Db) Select(query string) []oneRow {
	rows, err := db.connection.Query(query)

	checkError(err)
	defer rows.Close()

	columns, err := rows.Columns()

	checkError(err)


	rowsMap := make([]oneRow, 0)

	for rows.Next() {

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))

		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		checkError(err)

		oneRowMap := make(oneRow)

		for i, col := range values {
			oneRowMap[columns[i]] = string(col)
		}
		rowsMap = append(rowsMap, oneRowMap)
	}


	defer rows.Close()
	return rowsMap
}

func (db *Db)Close() {
	db.connection.Close()
}

func NewDb(driverName, dsn string) *Db {
	connection, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/jiajun")

	checkError(err)

	return &Db{
		dsn:      dsn,
		connection: connection,
	}
}


func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

