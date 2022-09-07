
-- book table
create table if not exists `t_books` (
  `f_isbn` varchar(14) not null,
  `f_title` varchar(200) default null,
  `f_price` int(11) default null,
  primary key (`f_isbn`)
) engine=innodb default charset=utf8;


create table if not exists t_countries (
  f_country_id varchar (2),
  f_country_name varchar (40),
  f_region_id decimal (10, 0)
);
