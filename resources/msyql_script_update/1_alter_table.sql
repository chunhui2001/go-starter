
-- add column
SELECT
	count(*) INTO @exist
FROM
	information_schema. COLUMNS
WHERE
	table_schema = 'mydb'
AND COLUMN_NAME = 'f_author'
AND table_name = 't_books'
LIMIT 1;


SET @QUERY =
IF (
	@exist <= 0,
	"ALTER TABLE `t_books` ADD `f_author` int(1) NOT NULL default '0'",
	'select \'Column Exists\' status'
);

PREPARE stmt
FROM
	@QUERY;

EXECUTE stmt;


-- add f_id
SELECT
	count(*) INTO @exist
FROM
	information_schema. COLUMNS
WHERE
	table_schema = 'mydb'
AND COLUMN_NAME = 'f_id'
AND table_name = 't_books'
LIMIT 1;


SET @QUERY =
IF (
	@exist <= 0,
	"ALTER TABLE `mydb`.`t_books` ADD COLUMN `f_id` int NOT NULL AUTO_INCREMENT FIRST, CHANGE COLUMN `f_isbn` `f_isbn` varchar(14) NOT NULL AFTER `f_id`, CHANGE COLUMN `f_title` `f_title` varchar(200) DEFAULT NULL AFTER `f_isbn`, CHANGE COLUMN `f_price` `f_price` int(11) DEFAULT NULL AFTER `f_title`, CHANGE COLUMN `f_author` `f_author` int(1) NOT NULL DEFAULT 0 AFTER `f_price`, DROP PRIMARY KEY, ADD PRIMARY KEY (`f_id`);",
	'select \'Column Exists\' status'
);

PREPARE stmt
FROM
	@QUERY;

EXECUTE stmt;



-- add f_created_at
SELECT
	count(*) INTO @exist
FROM
	information_schema. COLUMNS
WHERE
	table_schema = 'mydb'
AND COLUMN_NAME = 'f_created_at'
AND table_name = 't_books'
LIMIT 1;


SET @QUERY =
IF (
	@exist <= 0,
	"ALTER TABLE `mydb`.`t_books` ADD COLUMN `f_created_at` timestamp(3) NULL DEFAULT CURRENT_TIMESTAMP(3) AFTER `f_author`;",
	'select \'Column Exists\' status'
);

PREPARE stmt
FROM
	@QUERY;

EXECUTE stmt;
