

-- 查询某个表的所有字段的数据类型
SELECT distinct upper(DATA_TYPE) FROM INFORMATION_SCHEMA.COLUMNS 
  WHERE table_name = 't_token_price' ;

VARCHAR
BIGINT
LONGTEXT
DATETIME
INT
TINYINT
DECIMAL
DOUBLE
TIMESTAMP
CHAR
SET
ENUM
FLOAT
LONGBLOB
MEDIUMTEXT
MEDIUMBLOB
SMALLINT
TEXT
BLOB
TIME


-- 查询当前数据开启的所有事物
SELECT
  trx_mysql_thread_id,
  trx_id,
  trx_state,
  trx_started,
  trx_rows_locked,
  trx_query,
  trx_rows_locked,
  trx_isolation_level
FROM
  information_schema.innodb_trx;

-- 查询连接到当前服务器上的所有连接
SELECT
  tmp.ipAddress,
  -- Calculate how many connections are being held by this IP address.
  COUNT(*) AS numConnections,
  -- For each connection, the TIME column represent how many SECONDS it has been in
  -- its current state. Running some aggregates will give us a fuzzy picture of what
  -- the connections from this IP address is doing.
  FLOOR(AVG(tmp.time)) AS timeAVG,
  MAX(tmp.time) AS timeMAX
FROM
  -- Create an intermediary table that includes an additional column representing
  -- the client IP address without the port.
  (
    SELECT
      -- We don't actually need all of these columns but, including them here to
      -- demonstrate what fields COULD be used in the processlist system.
      pl.id,
      pl. USER,
      pl. HOST,
      pl.db,
      pl.command,
      pl.time,
      pl.state,
      pl.info,
      -- The host column is in the format of "IP:PORT". We want to strip off
      -- the port number so that we can group the results by the IP alone.
      LEFT (
        pl. HOST,
        (LOCATE(':', pl. HOST) - 1)
      ) AS ipAddress
    FROM
      INFORMATION_SCHEMA. PROCESSLIST pl
  ) AS tmp
GROUP BY
  tmp.ipAddress
ORDER BY
  numConnections DESC;



