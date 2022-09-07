
create table if not exists `t_token_list` (
  `f_id` int(11) unsigned not null auto_increment comment '自增主键',
  `f_chain_id` int(4) not null comment '链id',
  `f_token_name` varchar(65) collate utf8mb4_unicode_ci not null comment '代币名字',
  `f_symbol` varchar(625) collate utf8mb4_unicode_ci not null comment '代币符号',
  `f_token_type` int(2) default null comment '代币类型, 1: 是一个正常的代币, 2: 不是一个正常的代币',
  `f_addr` varchar(125) collate utf8mb4_unicode_ci not null comment '代币合约地址',
  `f_decimals` int(2) default null comment '精度',
  `f_supply` decimal(65,18) not null default '0.000000000000000000' comment '发行量',
  `f_icon` varchar(625) collate utf8mb4_unicode_ci default null comment '代币logo图标',
  `f_website` varchar(625) collate utf8mb4_unicode_ci default null comment '官方网站',
  `f_created_at` timestamp(3) null default null comment '创建时间',
  primary key (`f_id`),
  unique key `uniq_chain_addr` (`f_chain_id`,`f_addr`) using btree
) engine=innodb auto_increment=470 default charset=utf8mb4 collate=utf8mb4_unicode_ci comment='默认代币基础信息表';
