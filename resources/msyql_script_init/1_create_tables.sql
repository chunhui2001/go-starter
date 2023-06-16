
-- user table
CREATE TABLE if not exists t_user (
  f_id int(11) unsigned NOT NULL AUTO_INCREMENT,
  f_age INT,
  f_first_name varchar(200),
  f_last_name varchar(200),
  f_email varchar(200) UNIQUE NOT NULL,
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

-- 专辑表
CREATE TABLE if not exists t_album (
  f_id int(11) unsigned NOT NULL AUTO_INCREMENT,
  f_name varchar(200),
  f_title varchar(200),
  f_author varchar(200),
  f_quantity int(11),
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

-- 专辑订单表
CREATE TABLE if not exists t_album_order (
  f_id int(11) unsigned NOT NULL AUTO_INCREMENT,
  f_album_id int(11) unsigned NOT NULL, 
  f_cust_id int(11) unsigned NOT NULL, 
  f_quantity int(11), 
  f_created_at timestamp(3) not NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`f_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

