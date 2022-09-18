package actions

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	. "github.com/chunhui2001/go-starter/core/commons"
	"github.com/chunhui2001/go-starter/core/gsql"
	_ "github.com/chunhui2001/go-starter/core/utils"

	"github.com/gin-gonic/gin"
)

// MySQL发生死锁有哪些原因
// https://blog.csdn.net/weixin_42113222/article/details/115210334
// 查询当前数据库运行的所有事务
// select trx_mysql_thread_id,trx_id,trx_state,trx_started,trx_rows_locked,trx_query,trx_rows_locked,trx_isolation_level from information_schema.innodb_trx;
// 查询当前数据库运行的所有锁
// SELECT * FROM information_schema.innodb_locks
// 锁等待的对应关系，查看等待锁的事务
// SELECT * FROM information_schema.INNODB_LOCK_WAITS

func MySqlTrxLocks1(c *gin.Context) {

	fail := func(memo string, err error) error {
		return fmt.Errorf("Query-Album-Error: memo:%s, %v", memo, err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)

	defer cancelFunc()

	var rows *sql.Rows
	var err error
	var tx *sql.Tx

	// Get a Tx for making transaction requests.
	tx, err = gsql.DbClient.BeginTx(ctx, nil)

	if err != nil {
		c.JSON(200, (&R{Error: fail("BeginTxError", err)}).Fail(500))
		return
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	rows, err = tx.QueryContext(ctx, `select * from t_album where f_id in (10,8,5) for update;`)

	time.Sleep(1000 * time.Second)

	if err != nil {
		c.JSON(200, (&R{Error: fail("QueryContextError", err)}).Fail(500))
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		c.JSON(200, (&R{Error: fail("RowsColumnsError", err)}).Fail(500))
		return
	}

	colTypes, err2 := rows.ColumnTypes()

	if err2 != nil {
		c.JSON(200, (&R{Error: fail("ColumnTypesError", err)}).Fail(500))
		return
	}

	var result []map[string]interface{}

	for rows.Next() {

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
			c.JSON(200, (&R{Error: fail("RowsScanError", err)}).Fail(500))
			return
		}

		currentRow := make(map[string]interface{})

		for i, colName := range cols {
			currentRow[colName] = values[i]
		}

		result = append(result, currentRow)
	}

	if len(result) == 0 {
		result = []map[string]interface{}{}
	}

	// // Commit the transaction.
	// if err = tx.Commit(); err != nil {
	// 	c.JSON(200, (&R{Error: fail("CommitError", err)}).Fail(500))
	// 	return
	// }

	c.JSON(200, (&R{Data: result}).Success())

}

// https://blog.csdn.net/qq_45830276/article/details/125246751
// --------------------------
// 1、如何保证原子性：
// --------------------------
// 首先：对于A和B两操作要操作成功就一定需要更改到表的信息，如果如图所示A语句操作成功，而B语句操作时出现断电等其他情况终止了操作，
// 所以此时两个事务没有操作成功，在没有提交事务之前，mysql 会先记录更新前的数据到 undo log 日志里面，
// 当最终的因为操作不成功而发生事务回滚时，会从 undo log 日志里面先前存好的数据，重新对数据库的数据进行数据的回退。
// undo log 日志:（撤销回退的日志）主要存储数据库更新之前的数据，用于作备份

// --------------------------
// 2、如何保证事务的持久性：
// --------------------------
// 通过重做日志: redo log 日志，对于用户将对发生了修改而为提交的数据存入了 redo log 日志中，当此时发生断电等其他异常时，
// 可以根据 redo log 日志重新对数据做一个提交，做一个恢复。

// --------------------------
// 持久性产生的问题:
// --------------------------

// --------------------------
// 隔离性的隔离级别
// --------------------------
// 读未提交 read uncommitted
// 读已提交 read committed
// 可重复读 repeatable read
// 串行化 serializable

// (1)、读未提交：
// 事物A和事物B，事物A未提交的数据，事物B可以读取到。
// 这种隔离级别最低，这种级别一般是在理论上存在，数据库隔离级别一般都高于该级别。
// 三种并发问题都没解决。
// (2)、读已提交：
// 事务A只能读取到事务B提交的数据，这种级别可以避免“脏数据” ，这种隔离级别会导致“不可重复读取” ，Oracle默认隔离级别
// (3)、可重复读：(对于InnoDB不可能)
// 事务A和事务B，事务A提交之后的数据，事务B读取不到 - 事务B是可重复读取数据 - 这种隔离级别高于读已提交
// 换句话说，对方提交之后的数据，我还是读取不到 - 这种隔离级别可以避免“不可重复读取”，达到可重复读取 -
// 比如1点和2点读到数据是同一个 - MySQL默认级别 - 虽然可以达到可重复读取，但是会导致“幻像读”
// (4)、串行化：
// 事务A和事务B，事务A在操作数据库时，事务B只能排队等待 这种隔离级别很少使用，吞吐量太低，
// 用户体验差 这种级别可以避免“幻像读”，每一次读取的都是数据库中真实存在数据，事务A与事务B串行，而不并发
func MySqlTxnsRouter(c *gin.Context) {
	c.JSON(200, (&R{Data: true}).Success())
}
