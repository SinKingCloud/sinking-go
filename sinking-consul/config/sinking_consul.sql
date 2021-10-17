# Host: 42.157.128.40  (Version: 5.7.35)
# Date: 2021-10-17 21:31:11
# Generator: MySQL-Front 5.3  (Build 4.234)

/*!40101 SET NAMES utf8 */;

#
# Structure for table "sinking_apps"
#

CREATE TABLE `sinking_apps` (
  `id` bigint(3) NOT NULL DEFAULT '0',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '标识',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`),
  KEY `idx_all` (`title`,`name`,`create_time`,`update_time`,`is_delete`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应用列表';

#
# Data for table "sinking_apps"
#

INSERT INTO `sinking_apps` VALUES (1,'默认应用','default','2021-10-16 15:35:00','2021-10-16 15:35:00',0);

#
# Structure for table "sinking_configs"
#

CREATE TABLE `sinking_configs` (
  `id` bigint(3) NOT NULL DEFAULT '0',
  `app_name` varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '' COMMENT '应用标识',
  `env_name` varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '' COMMENT '环境标识',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '标识',
  `content` text CHARACTER SET latin1 COMMENT '配置内容',
  `hash` varchar(100) CHARACTER SET latin1 NOT NULL DEFAULT '' COMMENT '配置md5 hash',
  `group_name` varchar(255) CHARACTER SET latin1 NOT NULL DEFAULT '' COMMENT '分组名称',
  `type` varchar(50) CHARACTER SET latin1 NOT NULL DEFAULT '' COMMENT '类型(txt|json|yaml|properties)',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(0:已发布/1:未发布)',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`),
  KEY `idx_all` (`title`,`app_name`,`env_name`,`group_name`,`name`,`hash`,`type`,`status`,`create_time`,`update_time`,`is_delete`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置表';

#
# Data for table "sinking_configs"
#

INSERT INTO `sinking_configs` VALUES (1,'default','dev','测试配置','test','{\"host\":\"127.0.0.9\"}','3ca42f984623d0b4eb8430e38f701791','default','json',0,'2021-10-14 17:49:58','2021-10-14 17:49:58',0);

#
# Structure for table "sinking_envs"
#

CREATE TABLE `sinking_envs` (
  `id` bigint(3) NOT NULL DEFAULT '0' COMMENT '应用环境',
  `app_id` bigint(3) NOT NULL DEFAULT '0' COMMENT '应用ID',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '标识',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`),
  KEY `idx_all` (`app_id`,`title`,`name`,`create_time`,`update_time`,`is_delete`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应用环境';

#
# Data for table "sinking_envs"
#

INSERT INTO `sinking_envs` VALUES (1,1,'开发环境','dev','2021-10-16 15:35:00','2021-10-16 15:35:00',0),(2,1,'生产环境','release','2021-10-16 15:35:00','2021-10-16 15:35:00',0);

#
# Structure for table "sinking_roles"
#

CREATE TABLE `sinking_roles` (
  `id` bigint(3) NOT NULL DEFAULT '0' COMMENT '角色id',
  `name` varchar(255) DEFAULT NULL COMMENT '角色名称',
  `auths` text COMMENT '授权访问路径',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`),
  KEY `idx_all` (`name`,`auths`(10),`create_time`,`update_time`,`is_delete`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

#
# Data for table "sinking_roles"
#


#
# Structure for table "sinking_users"
#

CREATE TABLE `sinking_users` (
  `id` bigint(3) NOT NULL DEFAULT '0',
  `role_id` bigint(3) NOT NULL DEFAULT '0' COMMENT '角色ID',
  `user` varchar(255) NOT NULL DEFAULT '' COMMENT '账号',
  `pwd` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(255) DEFAULT NULL COMMENT '账户名称',
  `login_time` datetime DEFAULT NULL COMMENT '登陆时间',
  `login_ip` varchar(255) DEFAULT NULL COMMENT '登陆ip',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `is_delete` int(11) NOT NULL DEFAULT '0' COMMENT '是否删除(0:否/1:是)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

#
# Data for table "sinking_users"
#

INSERT INTO `sinking_users` VALUES (1,0,'admin','$2a$04$nK49S2fPoyZdw727GBvd6upZN5mC/MPJfpY3t9qsSmAOADdjhZuV.','管理员','2021-10-17 13:37:41','112.32.153.119',NULL,'2021-10-17 13:37:41',0);
