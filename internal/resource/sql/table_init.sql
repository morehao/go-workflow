CREATE TABLE `user`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `company_id`      bigint                                                       NOT NULL DEFAULT '0' COMMENT '公司id',
    `company_name`    varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司名称',
    `department_id`   bigint                                                       NOT NULL DEFAULT '0' COMMENT '部门id',
    `department_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '部门名称',
    `name`            varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE `procdef`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `name`        varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '流程名称',
    `resource`    json                                                                  DEFAULT NULL COMMENT '审批配置',
    `company`     varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司名称',
    `userid`      varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
    `username`    varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',
    `version`     int                                                          NOT NULL DEFAULT '0' COMMENT '流程版本',
    `deploy_time` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '部署时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='审批流程定义表';

CREATE TABLE `procdef_history`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `name`        varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '流程名称',
    `resource`    json                                                                  DEFAULT NULL COMMENT '审批配置',
    `company`     varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公司名称',
    `userid`      varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户id',
    `username`    varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',
    `version`     int                                                          NOT NULL DEFAULT '0' COMMENT '流程版本',
    `deploy_time` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '部署时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='审批流程定义历史记录表';

CREATE TABLE `proc_inst`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `proc_def_id`     bigint NOT NULL                                               DEFAULT 0 COMMENT '流程定义ID',
    `proc_def_name`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程定义名称',
    `title`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程实例标题',
    `department`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部门名称',
    `company`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '公司名称',
    `node_id`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '当前节点ID',
    `candidate`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '候选人',
    `task_id`         bigint NOT NULL                                               DEFAULT 0 COMMENT '任务ID',
    `start_time`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '开始时间',
    `end_time`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '结束时间',
    `duration`        bigint                                                        DEFAULT NULL COMMENT '持续时间',
    `start_user_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发起用户ID',
    `start_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发起用户名称',
    `is_finished`     tinyint                                                       DEFAULT '0' COMMENT '是否完成，1：已完成，2：未完成',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程实例表';

CREATE TABLE `proc_inst_history`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `proc_def_id`     bigint                                                        DEFAULT NULL COMMENT '流程定义ID',
    `proc_def_name`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程定义名称',
    `title`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '流程实例标题',
    `department`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部门名称',
    `company`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '公司名称',
    `node_id`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '当前节点ID',
    `candidate`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '候选人',
    `task_id`         bigint                                                        DEFAULT NULL COMMENT '任务ID',
    `start_time`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '开始时间',
    `end_time`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '结束时间',
    `duration`        bigint                                                        DEFAULT NULL COMMENT '持续时间',
    `start_user_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发起用户ID',
    `start_user_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发起用户名称',
    `is_finished`     tinyint(1) DEFAULT '0' COMMENT '是否完成',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='流程实例历史记录表';

CREATE TABLE `execution`
(
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `rev`           bigint                                                         DEFAULT NULL COMMENT '修订版本号',
    `proc_inst_id`  bigint                                                         DEFAULT NULL COMMENT '流程实例ID',
    `proc_def_id`   bigint                                                         DEFAULT NULL COMMENT '流程定义ID',
    `proc_def_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '流程定义名称',
    `node_infos`    varchar(4000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点信息',
    `is_active`     tinyint                                                        DEFAULT NULL COMMENT '是否活跃',
    `start_time`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '开始时间',
    PRIMARY KEY (`id`),
    KEY             `idx_proc_inst_id` (`proc_inst_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='执行实例表';

CREATE TABLE `execution_history`
(
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `rev`           int                                                            DEFAULT NULL COMMENT '修订版本号',
    `proc_inst_id`  bigint                                                         DEFAULT NULL COMMENT '流程实例ID',
    `proc_def_id`   bigint                                                         DEFAULT NULL COMMENT '流程定义ID',
    `proc_def_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '流程定义名称',
    `node_infos`    varchar(4000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点信息',
    `is_active`     tinyint                                                        DEFAULT NULL COMMENT '是否活跃',
    `start_time`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT NULL COMMENT '开始时间',
    PRIMARY KEY (`id`),
    KEY             `idx_proc_inst_id` (`proc_inst_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='执行实例历史记录表';

CREATE TABLE `identitylink`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '组',
    `type`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型',
    `user_id`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户ID',
    `user_name`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名称',
    `task_id`      bigint                                                        DEFAULT NULL COMMENT '任务ID',
    `step`         int                                                           DEFAULT NULL COMMENT '步骤',
    `proc_inst_id` bigint                                                        DEFAULT NULL COMMENT '流程实例ID',
    `company`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '公司',
    `comment`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '评论',
    PRIMARY KEY (`id`),
    KEY            `idx_proc_inst_id` (`proc_inst_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='身份链接表';

CREATE TABLE `identitylink_history`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `group`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '组',
    `type`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型',
    `user_id`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户ID',
    `user_name`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名称',
    `task_id`      bigint                                                        DEFAULT NULL COMMENT '任务ID',
    `step`         int                                                           DEFAULT NULL COMMENT '步骤',
    `proc_inst_id` bigint                                                        DEFAULT NULL COMMENT '流程实例ID',
    `company`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '公司',
    `comment`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '评论',
    PRIMARY KEY (`id`),
    KEY            `idx_proc_inst_id` (`proc_inst_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='身份链接历史记录表';

CREATE TABLE `task`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `node_id`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点ID',
    `step`            bigint                                                        DEFAULT NULL COMMENT '步骤',
    `proc_inst_id`    bigint                                                        DEFAULT NULL COMMENT '流程实例ID',
    `assignee`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '受派人',
    `create_time`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建时间',
    `claim_time`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '认领时间',
    `member_count`    tinyint                                                       DEFAULT '1' COMMENT '成员数量',
    `un_complete_num` tinyint                                                       DEFAULT '1' COMMENT '未完成数量',
    `agree_num`       tinyint                                                       DEFAULT NULL COMMENT '同意数量',
    `act_type`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'or' COMMENT '行为类型',
    `is_finished`     tinyint(1) DEFAULT '0' COMMENT '是否完成',
    PRIMARY KEY (`id`),
    KEY               `idx_proc_inst_id` (`proc_inst_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务表';

CREATE TABLE `task_history`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `node_id`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '节点ID',
    `step`            int                                                           DEFAULT NULL COMMENT '步骤',
    `proc_inst_id`    bigint                                                        DEFAULT NULL COMMENT '流程实例ID',
    `assignee`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '受派人',
    `create_time`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建时间',
    `claim_time`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '认领时间',
    `member_count`    tinyint                                                       DEFAULT '1' COMMENT '成员数量',
    `un_complete_num` tinyint                                                       DEFAULT '1' COMMENT '未完成数量',
    `agree_num`       tinyint                                                       DEFAULT NULL COMMENT '同意数量',
    `act_type`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'or' COMMENT '行为类型',
    `is_finished`     tinyint(1) DEFAULT '0' COMMENT '是否完成',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务历史记录表';
