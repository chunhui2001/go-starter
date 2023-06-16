package gsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
)

// ss := &SimpleSelect{
// 			Table:  "t_books",
// 			Fields: []string{"f_id", "f_title", "f_created_at"},
// 			Params: utils.MapOf("f_id", 1, "f_title", "sd"),
// 		}

func SimpleQueryWithContext(selects *SimpleSelect) ([]map[string]interface{}, error) {

	xselect, valus := selects.ToString()

	fail := func(err error) error {
		return fmt.Errorf("SimpleQueryWithContextError: %v", err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond*1000000)

	defer cancelFunc()

	var rows *sql.Rows
	var err error
	var tx *sql.Tx

	if selects.BeginTrx {

		// Get a Tx for making transaction requests.
		tx, err = DbClient.BeginTx(ctx, nil)

		if err != nil {
			return nil, fail(err)
		}

		// Defer a rollback in case anything fails.
		defer tx.Rollback()

		rows, err = tx.QueryContext(ctx, xselect, valus...)

	} else {
		rows, err = DbClient.QueryContext(ctx, xselect, valus...)
	}

	if err != nil {
		logger.Infof("Mysql-SimpleQuery-Error: sql=%s, valus=%s, error=%s", xselect, utils.ToJsonString(valus), err.Error())
		return nil, fail(err)
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if cols == nil {
		return nil, fail(err)
	}

	colTypes, err2 := rows.ColumnTypes()

	if err2 != nil {
		return nil, fail(err2)
	}

	var result []map[string]interface{}

	for rows.Next() {

		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		values := make([]interface{}, len(cols))

		for i := range values {
			currType := colTypes[i].DatabaseTypeName()
			if currType == "INT" {
				values[i] = new(int32)
			} else if currType == "VARCHAR" {
				values[i] = new(string)
			} else if currType == "TIMESTAMP" {
				values[i] = new(string)
			} else if currType == "TEXT" {
				values[i] = new(string)
			} else if currType == "UNSIGNED BIGINT" {
				values[i] = new(uint64)
			} else {
				logger.Errorf("Mysql-Current-Data-Type-Not-Cached: DatabaseTypeName=%s, ColumeName=%s", currType, colTypes[i].Name())
				values[i] = new(interface{})
			}
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(values...); err != nil {
			return nil, fail(err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		currentRow := make(map[string]interface{})

		for i, colName := range cols {
			currentRow[colName] = values[i]
		}

		result = append(result, currentRow)
	}

	if selects.BeginTrx {
		// Commit the transaction.
		if err = tx.Commit(); err != nil {
			return nil, fail(err)
		}
	}

	return result, nil

}

func SimpleQuery(selects *SimpleSelect) ([]map[string]interface{}, error) {

	xselect, valus := selects.ToString()
	rows, err := DbClient.Query(xselect, valus...)

	if err != nil {
		logger.Infof("Mysql-SimpleQuery-Error: sql=%s, valus=%s, error=%s", xselect, utils.ToJsonString(valus), err.Error())
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if cols == nil {
		return nil, err
	}

	colTypes, err2 := rows.ColumnTypes()

	if err2 != nil {
		return nil, err2
	}

	var result []map[string]interface{}

	for rows.Next() {

		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		values := make([]interface{}, len(cols))

		for i := range values {
			currType := colTypes[i].DatabaseTypeName()
			if currType == "INT" {
				values[i] = new(int32)
			} else if currType == "VARCHAR" {
				values[i] = new(string)
			} else if currType == "TIMESTAMP" {
				values[i] = new(string)
			} else {
				logger.Errorf("Mysql-Current-Data-Type-Not-Cached: DatabaseTypeName=%s, ColumeName=%s", currType, colTypes[i].Name())
				values[i] = new(interface{})
			}
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		currentRow := make(map[string]interface{})

		for i, colName := range cols {
			currentRow[colName] = values[i]
		}

		result = append(result, currentRow)
	}

	return result, nil

}
