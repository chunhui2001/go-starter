
-- add column
SELECT
	count(*) INTO @exist
FROM
	information_schema. COLUMNS
WHERE
	table_schema = 'mydb'
AND table_name = 't_user'
AND COLUMN_NAME = 'f_created_at'
LIMIT 1;


SET @QUERY =
IF (
	@exist <= 0,
	"ALTER TABLE `t_user` ADD `f_created_at` DATETIME(3) not NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'",
	'select \'Column Exists\' status'
);

PREPARE stmt
FROM
	@QUERY;

EXECUTE stmt;


-- add column
SELECT
	count(*) INTO @exist
FROM
	information_schema. COLUMNS
WHERE
	table_schema = 'mydb'
AND table_name = 't_album'
AND COLUMN_NAME = 'f_created_at'
LIMIT 1;


SET @QUERY =
IF (
	@exist <= 0,
	"ALTER TABLE `t_album` ADD `f_created_at` DATETIME(3) not NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间'",
	'select \'Column Exists\' status'
);

PREPARE stmt
FROM
	@QUERY;

EXECUTE stmt;
