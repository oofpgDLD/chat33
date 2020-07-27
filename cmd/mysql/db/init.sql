/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1（docker）
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 127.0.0.1:3306
 Source Schema         : chat33

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : 65001

 Date: 22/07/2020 10:19:40
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

grant all PRIVILEGES on *.* to root@'%' identified by '123456';
flush privileges;

-- ----------------------------
-- create database
-- ----------------------------
create database `chat33`; 
SET character_set_client = utf8mb4; 
use chat33;

-- ----------------------------
-- Table structure for add_friend_conf
-- ----------------------------
DROP TABLE IF EXISTS `add_friend_conf`;
CREATE TABLE `add_friend_conf`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `need_confirm` int(2) NOT NULL COMMENT '1 需要验证 2  不需要验证',
  `need_answer` int(2) NOT NULL COMMENT '1 需要回答问题 2  不需要回答',
  `question` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '问题',
  `answer` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '答案',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin`  (
  `id` int(11) NOT NULL,
  `app_id` int(11) NULL DEFAULT NULL,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '账号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
  `salt` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '盐值',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `app_id_account_UNIQUE`(`app_id`, `account`) USING BTREE COMMENT 'app_id和account唯一键'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for admin_operate_log
-- ----------------------------
DROP TABLE IF EXISTS `admin_operate_log`;
CREATE TABLE `admin_operate_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `operator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '管理员id',
  `type` tinyint(2) NULL DEFAULT NULL COMMENT '1 群 2用户',
  `target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '操作对象id',
  `operate_type` tinyint(2) NULL DEFAULT NULL COMMENT '操作类型：1 封号 2 封群 3 解封号 4解封群',
  `reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `effective_time` bigint(20) NULL DEFAULT NULL COMMENT '截止有效时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for advertisement
-- ----------------------------
DROP TABLE IF EXISTS `advertisement`;
CREATE TABLE `advertisement`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '广告id',
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '广告名称',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '广告地址',
  `duration` int(255) NULL DEFAULT 3 COMMENT '持续时长 单位：s秒',
  `link` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '跳转地址',
  `is_active` int(255) NULL DEFAULT 1 COMMENT '是否激活 0：未激活 1：激活',
  `is_delete` int(255) NULL DEFAULT 0 COMMENT '是否删除 0：未删除 1：删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app`  (
  `app_id` int(5) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '名称',
  `user_info_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `redpacket_pid` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '红包对应平台id',
  `redpacket_server` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '红包服务地址',
  `redpacket_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '红包h5地址',
  `is_inner` tinyint(1) NULL DEFAULT 0 COMMENT '是否是内部 账户体系（托管账户）',
  `main_coin` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '主要币种',
  `backend_app_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `backend_app_secret` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `push_app_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '友盟推送 appKey',
  `push_app_master_secret` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '友盟推送 app_master_secret',
  `push_mi_active` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '友盟推送 系统推送时打开的active',
  `is_otc` tinyint(255) NULL DEFAULT NULL COMMENT '是否是otc模式',
  `otc_server` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'otc服务访问地址',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`app_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for app_module
-- ----------------------------
DROP TABLE IF EXISTS `app_module`;
CREATE TABLE `app_module`  (
  `app_id` int(11) NOT NULL COMMENT '应用类型',
  `type` tinyint(255) NOT NULL COMMENT '模块类型',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '模块名称',
  `enable` tinyint(255) NULL DEFAULT 0 COMMENT '是否使能 0不使能 1使能',
  PRIMARY KEY (`app_id`, `type`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for app_update
-- ----------------------------
DROP TABLE IF EXISTS `app_update`;
CREATE TABLE `app_update`  (
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型昵称',
  `compatible` tinyint(255) NULL DEFAULT NULL COMMENT '是否开启兼容模式：0 否 1 是',
  `min_version_code` int(255) NULL DEFAULT NULL COMMENT '兼容的最小code  iso计算是版本号第一位*10000 + 第二位*100 + 第三位',
  `version_code` int(255) NULL DEFAULT NULL COMMENT '版本code',
  `version_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '版本号名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '更新内容描述',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '更新包地址',
  `size` bigint(255) NULL DEFAULT NULL COMMENT '更新包大小',
  `md5` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `force` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`app_id`, `name`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for apply
-- ----------------------------
DROP TABLE IF EXISTS `apply`;
CREATE TABLE `apply`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` tinyint(255) NULL DEFAULT NULL COMMENT '1群 2好友',
  `apply_user` int(255) NULL DEFAULT NULL,
  `target` int(255) NULL DEFAULT NULL,
  `apply_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `state` tinyint(255) NULL DEFAULT NULL COMMENT '1待处理   2拒绝   3同意',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `datetime` bigint(20) NULL DEFAULT NULL,
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '例 {\"sourceType\":1,\"sourceId\":\"123\"}',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_type_apply_user_target`(`type`, `apply_user`, `target`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for client_log
-- ----------------------------
DROP TABLE IF EXISTS `client_log`;
CREATE TABLE `client_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` int(11) NOT NULL COMMENT '1 连接 2断开',
  `user_id` int(11) NULL DEFAULT NULL,
  `client_uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `time` bigint(20) NULL DEFAULT NULL,
  `pushed` tinyint(2) NULL DEFAULT 1 COMMENT '1:未推送 2: 已推送',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for coin
-- ----------------------------
DROP TABLE IF EXISTS `coin`;
CREATE TABLE `coin`  (
  `coin_id` int(255) NOT NULL,
  `coin_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`coin_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for friends
-- ----------------------------
DROP TABLE IF EXISTS `friends`;
CREATE TABLE `friends`  (
  `user_id` int(11) NOT NULL,
  `friend_id` int(11) NOT NULL,
  `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '好友备注',
  `add_time` bigint(20) NOT NULL DEFAULT 0,
  `DND` tinyint(4) NOT NULL DEFAULT 2 COMMENT '是否消息免打扰 1免打扰  2关闭',
  `top` int(3) NOT NULL DEFAULT 2 COMMENT '好友置顶  1置顶 2不置顶',
  `type` int(3) NULL DEFAULT 1 COMMENT '1 普通  2 常用',
  `is_blocked` int(3) NULL DEFAULT 0 COMMENT '0 非黑名单 1黑名单',
  `is_delete` int(3) NOT NULL DEFAULT 1 COMMENT '1 未删除  2 已删除',
  `source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `ext_remark` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '额外的备注信息',
  PRIMARY KEY (`user_id`, `friend_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for invite_room_conf
-- ----------------------------
DROP TABLE IF EXISTS `invite_room_conf`;
CREATE TABLE `invite_room_conf`  (
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户Id',
  `need_confirm` tinyint(255) NULL DEFAULT 0 COMMENT '0 不需要 1需要',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for login_log
-- ----------------------------
DROP TABLE IF EXISTS `login_log`;
CREATE TABLE `login_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `login_time` bigint(20) NOT NULL,
  `device` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '登录设备',
  `device_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `login_type` tinyint(3) NULL DEFAULT NULL COMMENT '1手机验证，2手机密码，3邮箱验证，4邮箱密码',
  `client_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `uuid` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `version` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT 'app当前版本号',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for open_log
-- ----------------------------
DROP TABLE IF EXISTS `open_log`;
CREATE TABLE `open_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `device` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'app版本号',
  `create_time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for praise
-- ----------------------------
DROP TABLE IF EXISTS `praise`;
CREATE TABLE `praise`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `channel_type` tinyint(2) NULL DEFAULT NULL COMMENT '1: 聊天室（弃用）； 2：群组；3：好友',
  `target_id` int(11) NULL DEFAULT NULL COMMENT '群或者好友id',
  `log_id` int(11) NULL DEFAULT NULL COMMENT '消息log',
  `sender_id` int(11) NULL DEFAULT NULL COMMENT '消息发出者id',
  `opt_id` int(11) NULL DEFAULT NULL COMMENT '赞赏者id',
  `type` tinyint(2) NULL DEFAULT NULL COMMENT '赞赏类型：1. 点赞；2.赏赐',
  `record_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '转账记录id',
  `coin_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '币种id',
  `coin_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '币种名称',
  `amount` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '金额',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '记录生成时间',
  `is_delete` tinyint(255) NULL DEFAULT NULL COMMENT '是否删除 0未删除 1 删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for praise_user
-- ----------------------------
DROP TABLE IF EXISTS `praise_user`;
CREATE TABLE `praise_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `target_id` int(11) NULL DEFAULT NULL COMMENT '好友id',
  `opt_id` int(11) NULL DEFAULT NULL COMMENT '赞赏者id',
  `record_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '转账记录id',
  `coin_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '币种id',
  `coin_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '币种名称',
  `amount` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '金额',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '记录生成时间',
  `is_delete` tinyint(255) NULL DEFAULT NULL COMMENT '是否删除 0未删除 1 删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for private_chat_log
-- ----------------------------
DROP TABLE IF EXISTS `private_chat_log`;
CREATE TABLE `private_chat_log`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `msg_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '消息id,客户端生成',
  `sender_id` int(11) NULL DEFAULT NULL,
  `receive_id` int(11) NULL DEFAULT NULL,
  `is_snap` tinyint(2) NULL DEFAULT 2 COMMENT '1:是阅后即焚消息 2：不是',
  `msg_type` int(10) UNSIGNED NULL DEFAULT NULL,
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `status` int(11) NULL DEFAULT NULL COMMENT '1 已读  2未读',
  `send_time` bigint(13) NULL DEFAULT NULL,
  `ext` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `is_delete` int(3) NOT NULL DEFAULT 2 COMMENT '1删除  2未删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `sender_id`(`sender_id`) USING BTREE,
  INDEX `receive_id`(`receive_id`) USING BTREE,
  INDEX `send_time`(`send_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for push
-- ----------------------------
DROP TABLE IF EXISTS `push`;
CREATE TABLE `push`  (
  `device_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '推动设备识别码',
  `user_id` int(10) NULL DEFAULT NULL COMMENT '用户id',
  `device_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '设备类型',
  PRIMARY KEY (`device_token`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for red_packet_log
-- ----------------------------
DROP TABLE IF EXISTS `red_packet_log`;
CREATE TABLE `red_packet_log`  (
  `packet_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `ctype` tinyint(255) NOT NULL COMMENT '红包发到群/用户/聊天室',
  `user_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '发红包的用户id',
  `to_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `coin` int(11) NOT NULL,
  `size` int(11) NOT NULL COMMENT '几个红包',
  `amount` decimal(11, 0) NOT NULL COMMENT '总金额',
  `remark` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '红包备注',
  `type` int(11) NOT NULL COMMENT '新人、拼手气红包',
  `created_at` bigint(13) NOT NULL COMMENT '时间',
  PRIMARY KEY (`packet_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '发红包的记录' ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for room
-- ----------------------------
DROP TABLE IF EXISTS `room`;
CREATE TABLE `room`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mark_id` int(11) NULL DEFAULT NULL,
  `name` varchar(180) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `master_id` int(11) NULL DEFAULT NULL,
  `create_time` bigint(20) NULL DEFAULT NULL,
  `can_add_friend` tinyint(1) NULL DEFAULT 1 COMMENT '1可添加好友 2不可添加好友',
  `join_permission` tinyint(1) NULL DEFAULT 2 COMMENT '1需要审批 2不需要审批  3禁止加入',
  `record_permission` tinyint(1) NULL DEFAULT 1 COMMENT '1 可查看入群前记录 2 只可查看入群后记录',
  `admin_muted` tinyint(1) NULL DEFAULT NULL,
  `master_muted` tinyint(255) NULL DEFAULT NULL,
  `encrypt` tinyint(255) NULL DEFAULT 2 COMMENT '1 加密群 2非加密群',
  `is_delete` tinyint(255) NULL DEFAULT 1 COMMENT '1 未删除 2 删除',
  `room_level` int(255) NULL DEFAULT 1 COMMENT '群等级',
  `recommend` tinyint(2) NULL DEFAULT 0 COMMENT '0 非推荐群 1推荐群',
  `close_until` bigint(255) NULL DEFAULT 0 COMMENT '封群截止时间',
  `identification` tinyint(255) NULL DEFAULT NULL COMMENT '加v认证： 0 未认证 1已认证',
  `identification_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '加v认证描述信息',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `mark_id`(`mark_id`) USING BTREE,
  INDEX `master_id`(`master_id`) USING BTREE,
  INDEX `name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for room_config
-- ----------------------------
DROP TABLE IF EXISTS `room_config`;
CREATE TABLE `room_config`  (
  `level` tinyint(255) NOT NULL DEFAULT 1,
  `app_id` int(11) NOT NULL,
  `number_limit` int(255) NULL DEFAULT NULL COMMENT '群人数限制',
  PRIMARY KEY (`level`, `app_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for room_msg_content
-- ----------------------------
DROP TABLE IF EXISTS `room_msg_content`;
CREATE TABLE `room_msg_content`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `msg_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '消息id,客户端生成',
  `room_id` int(11) NULL DEFAULT NULL,
  `sender_id` int(11) NULL DEFAULT NULL,
  `is_snap` tinyint(2) NULL DEFAULT 2 COMMENT '1: ???ĺ????Ϣ 2?????',
  `msg_type` int(255) NULL DEFAULT NULL COMMENT '0：系统消息，1:文字，2:音频，3：图片，4：红包，5：视频',
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `datetime` bigint(20) NULL DEFAULT NULL,
  `ext` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `is_delete` tinyint(2) NULL DEFAULT 1 COMMENT '1 未删除 2 已删除',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `room_id`(`room_id`) USING BTREE,
  INDEX `sender_id`(`sender_id`) USING BTREE,
  INDEX `datetime`(`datetime`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for room_msg_receive
-- ----------------------------
DROP TABLE IF EXISTS `room_msg_receive`;
CREATE TABLE `room_msg_receive`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_msg_id` int(11) NOT NULL,
  `receive_id` int(11) NOT NULL,
  `state` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_room_msg_id_receive_id`(`room_msg_id`, `receive_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for room_user
-- ----------------------------
DROP TABLE IF EXISTS `room_user`;
CREATE TABLE `room_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `user_nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `level` int(11) NULL DEFAULT NULL COMMENT '1:普通用户 2：管理员 3：群主',
  `no_disturbing` tinyint(1) NULL DEFAULT NULL COMMENT '1：开启了免打扰，2：关闭',
  `common_use` tinyint(255) NULL DEFAULT NULL COMMENT '1 普通 ；2 常用',
  `room_top` tinyint(255) NULL DEFAULT NULL COMMENT '1 置顶 2不置顶 ',
  `create_time` bigint(20) NULL DEFAULT NULL,
  `is_delete` tinyint(255) NULL DEFAULT NULL COMMENT '1 未删除 2 删除',
  `source` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '例 {\"sourceType\":1,\"sourceId\":\"123\"}',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `room_id_user_id`(`room_id`, `user_id`) USING BTREE,
  INDEX `room_id`(`room_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for room_user_muted
-- ----------------------------
DROP TABLE IF EXISTS `room_user_muted`;
CREATE TABLE `room_user_muted`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) NULL DEFAULT NULL,
  `user_id` int(11) NULL DEFAULT NULL,
  `list_type` tinyint(2) NULL DEFAULT 1 COMMENT '1 未使用 2 黑名单 3 白名单',
  `deadline` bigint(255) NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `unique_room_id_user_id`(`room_id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS `token`;
CREATE TABLE `token`  (
  `user_id` int(255) NOT NULL COMMENT 'userid',
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'token',
  `time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `user_id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `mark_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '号码',
  `uid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用uid',
  `app_id` int(11) NULL DEFAULT NULL COMMENT '不用',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '手机号或者邮箱  找币那边的',
  `user_level` int(11) NULL DEFAULT NULL COMMENT '0 游客 1 普通用户 2 客服 3 管理员  数据库里不存游客  然后接口返回的话 客服和管理员都是2',
  `verified` int(11) NULL DEFAULT 0 COMMENT '是否实名  找币',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '管理员给的备注',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `sex` tinyint(4) NULL DEFAULT NULL,
  `area` varchar(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '区号',
  `phone` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `email` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `com_id` int(11) NULL DEFAULT NULL COMMENT '公司id 未用到',
  `device_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '友盟推送',
  `position` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '职位',
  `invite_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '邀请码',
  `deposit_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '钱包地址',
  `device` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `reg_version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '第一次注册版本',
  `now_version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '最新登录版本',
  `close_until` bigint(255) NULL DEFAULT 0 COMMENT '封号截止时间',
  `super_user_level` int(255) NULL DEFAULT 1 COMMENT '会员等级',
  `public_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '加密公钥',
  `private_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '私钥',
  `identification_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '加v认证描述信息',
  `identification` tinyint(255) NULL DEFAULT NULL COMMENT '加v认证： 0 未认证 1已认证',
  `ischain` tinyint(255) NULL DEFAULT 0 COMMENT '默认0未上链    0:未上链；1：已上链',
  PRIMARY KEY (`user_id`) USING BTREE,
  INDEX `create_time`(`create_time`) USING BTREE,
  INDEX `uid_UNIQUE`(`uid`, `app_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for user_config
-- ----------------------------
DROP TABLE IF EXISTS `user_config`;
CREATE TABLE `user_config`  (
  `level` tinyint(255) NOT NULL,
  `app_id` int(11) NOT NULL,
  `number_limit` int(255) NULL DEFAULT NULL COMMENT '创建群个数限制',
  PRIMARY KEY (`level`, `app_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for verify_apply
-- ----------------------------
DROP TABLE IF EXISTS `verify_apply`;
CREATE TABLE `verify_apply`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` tinyint(255) NULL DEFAULT NULL COMMENT '1个人认证；2群认证',
  `target_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '申请对象',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '认证描述',
  `amount` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '申请费用金额',
  `currency` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '申请费用币种',
  `state` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '审核状态：1：待审核；2已认证；3未通过',
  `record_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '转账记录id',
  `update_time` bigint(20) NULL DEFAULT NULL COMMENT '最近更新时间',
  `fee_state` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '0' COMMENT '手续费入账状态：0入账中 1入账成功 2入账失败 3退回中 4退回成功 5退回失败',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for verify_fee
-- ----------------------------
DROP TABLE IF EXISTS `verify_fee`;
CREATE TABLE `verify_fee`  (
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `type` tinyint(255) NOT NULL COMMENT '1：个人认证；2：群认证',
  `currency` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '币种',
  `amount` double(255, 0) NULL DEFAULT NULL COMMENT '金额',
  PRIMARY KEY (`app_id`, `type`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Triggers structure for table room
-- ----------------------------
DROP TRIGGER IF EXISTS `master_muted_type_update`;
delimiter ;;
CREATE TRIGGER `master_muted_type_update` AFTER UPDATE ON `room` FOR EACH ROW BEGIN
	if (NEW.master_muted != OLD.master_muted)
	THEN
	UPDATE room_user_muted 
	SET room_user_muted.list_type = 1 , room_user_muted.deadline = 0
	WHERE room_user_muted.room_id = NEW.id;
	END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO `chat33`.`admin`(`id`, `app_id`, `account`, `password`, `salt`) VALUES (1, 1001, 'admin', '315963ecbf53bb16fa186cc8e91225e91c6fe057911cc4f5ad2fdb7b720a8022', '6395ebd0f4b478145ecfbaf939454fa4');
INSERT INTO `chat33`.`app`(`app_id`, `name`, `user_info_url`, `redpacket_pid`, `redpacket_server`, `redpacket_url`, `is_inner`, `main_coin`, `backend_app_key`, `backend_app_secret`, `push_app_key`, `push_app_master_secret`, `push_mi_active`, `is_otc`, `otc_server`, `remark`) VALUES (1001, 'chat33Pro', '', '1005', '', '', 1, 'BTY', '', '', '{\r\n\"Android\": \"\",\r\n \"iOS\": \"\"\r\n}', '{\r\n\"Android\": \"\",\r\n \"iOS\": \"\"\r\n}', '{\r\n\"Android\": \"\",\r\n \"iOS\": \"\"\r\n}', NULL, '', '');
