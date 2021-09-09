# 实例管理表
create table `tb_walnut_data_source` (
  id int(11) unsigned not null auto_increment comment "主键id",
  db_addr varchar(64) not null comment "实例地址",
  db_account varchar(32) not null comment "实例帐号",
  db_password varchar(32) not null comment "实例密码",
  charset varchar() not null varchar(16) default "utf8mb4" comment "字符集",
  db_type tinyint(4) unsigned not null default 1 comment "数据库类型: 1-MySQL、2-Redis、3-MongoDB",
  input_mode tinyint(4) unsigned not null default 1 comment "实例录入方式: 1-ATOMIC、2-MHA、3-CDB",
  db_status tinyint(4) unsigned not null default 1 comment "实例状态: 0-关闭、1-开启",
  env tinyint(4) unsigned not null comment "环境类型: 1-prod、2-gray、3-test、4-dev",
  db_name varchar(32) not null comment "实例别名",
  belong_dba varchar(24) not null comment "实例所属DBA",
  change_table_no_lock tinyint(4) not null default 0 comment "不锁表变更: 0-关闭，1-开启",
  query_timeout tinyint(4) unsigned not null default 0 comment "查询超时时间, 0-表示不限制",
  export_timeout tinyint(4) unsigned not null default 0 comment "导出超时时间, 0-表示不限制",
  gmt_create int(11) not null default 0 comment "创建时间",
  gmt_modified int(11) not null default 0  comment "创建时间",
primary key (id)
) engine = innodb default character set = utf8mb4 comment "实例表(实例录入)";



CREATE TABLE `tb_walnut_apply_db` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `sponsor` varchar(24) NOT NULL COMMENT '发起人',
  `db_name` varchar(24) NOT NULL COMMENT '数据库名称',
  `charset` varchar(16) NOT NULL DEFAULT 'utf8mb4' COMMENT '字符集',
  `db_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '数据库类型: 1-MySQL、2-Redis、3-MongoDB',
  `env` tinyint(3) unsigned NOT NULL COMMENT '环境类型: 1-prod、2-gray、3-test、4-dev',
  `data_source_id` int(10) unsigned NOT NULL COMMENT '实例id, 来源tb_walnut_data_source.id',
  `source_name` varchar(32) not null comment "实例别名(实例来源名称)",
  `audit_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1-待审核,2-待执行,3-已完成,4-已失败,5-已拒绝,6-已撤销',
  `ext` varchar(128) DEFAULT NULL COMMENT '备注（数据库备注）',
  `audit_ext` varchar(128) DEFAULT NULL COMMENT '审核备注',
  `revoke_ext` varchar(128) DEFAULT NULL COMMENT '撤销备注',
  `gmt_create` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `gmt_modified` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_data_source_id` (`data_source_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源申请表';


# 工单管理
create table `tb_core_sql_order` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `order_id` varchar(32) not null comment "工单编号",
  `title` varchar(64) not null default "" COMMENT "工单标题",
  `demand` varchar(128) not null default "" COMMENT "需求背景",
  `sql_type` tinyint not null default 1 comment "工单类型，1-数据变更，2-结构变更，3-权限申请，4-数据导出",
  `sponsor` varchar(12) NOT NULL COMMENT '发起人',
  `reviewer` varchar(24) NOT NULL default "" COMMENT 'reviewer',
  `executor` varchar(24) NOT NULL default "" COMMENT 'executor',
  `env` tinyint(4) unsigned not null comment "环境类型: 1-生产、2-预发、3-测试、4-开发",
  `delay` varchar(24) not null default "" comment "定时执行，空表示不定时执行",
  `is_backup` tinyint(4) not null  default 1 COMMENT "0-否，1-是",
  `data_source_id` int not null COMMENT "实例id, 来源于tb_walnut_data_source.id",
  `remote_host` varchar(24) not null COMMENT "实例地址",
  `remote_port` int not null COMMENT "实例端口",
  `db_name` varchar(24) not null COMMENT "数据库名称",
  `progress` tinyint not null default 1 comment "工单进度，1-待审核，2-已驳回，3-待执行，4-执行中，5-已完成，6-已失败，7-已关闭",
  `step` tinyint not null default 1 COMMENT "当前审核步骤",
  `audit_step` tinyint not null default 1 comment "审核进度，多级审核，1-表示提交工单者，之后为审核者",
  `remark` varchar(128) not null default "" COMMENT "审核备注",
  `perm_detail` json COMMENT "权限申请详情",
  `contents` text COMMENT "提交的内容",
  `file_format` char(4) not null default 'xlsx' COMMENT "文件导出类型",
  `gmt_create` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `gmt_modified` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY(`id`),
  unique key `uk_orderid` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库申请工单';

alter table tb_core_sql_order add `perm_detail` json COMMENT "权限申请详情" after `remark`,


create table tb_core_sql_order_detail (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `order_id` varchar(32) not null comment "工单编号",
  `execute_result` json comment "执行结果",
  `gmt_create` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `gmt_modified` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间'
  PRIMARY KEY(`id`),
  unique key `uk_orderid` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工单详情';


create table `tb_core_work_flow` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `order_id` varchar(32) NOT NULL COMMENT "工单表id,来源tb_core_orders.order_id",
  `current_handler` varchar(12) not null COMMENT "当前处理人",
  `progress` tinyint(4) not null default 1 comment "工单进度，1-待审核，2-已驳回，3-已批准，4-处理中，5-已完成，6-已失败，7-已关闭",
  `contents` varchar(256) not null default "工单流转记录",
  `operate_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  `gmt_create` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `gmt_modified` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
 PRIMARY KEY(`id`),
 KEY idx_orderid(`order_id`, `operate_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工单流转记录';




# 权限管理
create table `tb_core_privileges` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `username` varchar(16) NOT NULL COMMENT "用户名称, 来源tb_walnut_account.username",
  `apply_type` tinyint(4) unsigned not null comment "权限列表: 1-schema, 2-table",
  `env` tinyint(4) unsigned not null comment "环境类型: 1-生产、2-预发、3-测试、4-开发",
  `perm_type` varchar(16) not null comment "权限类型，'1,2,3'(1-数据变更，2-结构变更，3-数据导出)",
  `data_source_id` int not null COMMENT "实例id, 来源于tb_walnut_data_source.id",
  `source_name` varchar(64) not null COMMENT "实例别名",
  `remote_host` varchar(24) not null COMMENT "实例地址",
  `remote_port` int not null COMMENT "实例端口",
  `schema` varchar(36) not null  COMMENT "数据库",
  `table` varchar(36) not null default "" COMMENT "数据表",
  `expire_time` int unsigned not null COMMENT "失效时间",
  `gmt_create` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `gmt_modified` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
 PRIMARY KEY(`id`),
 UNIQUE KEY uk_ua_dsid_st(`username`, `apply_type`, `data_source_id`, `schema`, `table`),
 KEY idx_username_expire(`username`,`env`,`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权限管理表';