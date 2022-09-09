package gsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Insert(`insert into book(isbn, title, price) values(?, ?, ?)`, "978-4798161495", "MySQL徹底入門 第4版", 4180)
func Insert(sql string, args ...any) sql.Result {

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Nanosecond*1000000)
	defer cancelFunc()

	result, err := DbClient.ExecContext(ctx, sql, args...)

	if err != nil {
		logger.Error(fmt.Sprintf("Mysql-Insert-Error: sql=%s, errorMessage=%s", sql, string(err.Error())))
		return nil
	}

	return result

}
