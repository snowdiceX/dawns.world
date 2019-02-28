/*
 * twilight RESTful api service
 *
 * login_user (api login)
 * -- 登录账户，不保存隐私信息
 * -- 功能：
 *     登录、登出
 *
 * wallet_address (api wallet)
 * -- 钱包地址，可多个网络，多种token，一个登录账户可多个钱包地址
 * -- 功能：
 *     充值、提现
 *
 * quotation (api quotation)
 * -- token 牌价
 *     牌价查询；token 兑换
 *
 * game (api game)
 * -- 游戏，用于记录游戏合约（联盟链合约），用户用来查看游戏项目和选择游戏
 * -- 功能：
 *     列表展示游戏类型，选择并进入游戏
 *
 * gaming (api gaming)
 * -- 正在进行中的游戏，处理游戏回合、游戏规则等游戏逻辑
 * -- 功能：
 *     游戏回合的处理，通用游戏处理，所有不同游戏类型的统一处理功能；
 *     进入游戏自动兑换筹码；游戏结束，自动筹码回兑、自动分配token
 *
 * public_random (api radom)
 * -- 公共随机共识（人肉共识），随机数共识生成合约（联盟链合约）
 * -- 功能：
 *     接收随机共识请求，调用合约生成随机数
 *
 * texas_poker (api texas_poker)
 * -- 德州扑克
 * -- 功能：
 *     游戏规则相关api
 *
 */

create database twilight;

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for login_user
-- ----------------------------
DROP TABLE IF EXISTS `login_user`;
CREATE TABLE `login_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'user id',
  `nick_name` varchar(64) DEFAULT NULL,
  `phone` varchar(128) DEFAULT '',
  `email` varchar(128) NOT NULL,
  `password` varchar(32) NOT NULL,
  `invite_code` varchar(10) DEFAULT '',
  `status` int(4) DEFAULT '0',
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `remarks` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=422 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for wallet_address
-- ----------------------------
DROP TABLE IF EXISTS `wallet_address`;
CREATE TABLE `wallet_address` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'wallet id',
  `user_id` int(11) NOT NULL,
  `network` varchar(32) DEFAULT NULL,
  `wallet_address` varchar(255) NOT NULL,
  `token_name` varchar(32) DEFAULT NULL,
  `balance` int(64) DEFAULT '0',
  `status` int(4) DEFAULT '0',
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `remarks` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_wallet_token` (`user_id`,`wallet_address`,`token_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for quotation
-- ----------------------------
DROP TABLE IF EXISTS `quotation`;
CREATE TABLE `quotation` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'token quotation id',
  `base_token_name` varchar(32) DEFAULT NULL,
  `network` varchar(32) DEFAULT NULL,
  `token_name` varchar(32) DEFAULT NULL,
  `price` int(11) DEFAULT '0',
  `status` int(4) DEFAULT '0',
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `remarks` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token_base` (`token_name`,`network`,`base_token_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for game
-- ----------------------------
DROP TABLE IF EXISTS `game`;
CREATE TABLE `game` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'game id',
  `contract_id` varchar(32) DEFAULT NULL COMMENT 'chaincode id',
  `title` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL,
  `status` int(4) DEFAULT '0',
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `remarks` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `contract_id` (`contract_id`)
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for gaming
-- ----------------------------
DROP TABLE IF EXISTS `gaming`;
CREATE TABLE `gaming` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'gaming id',
  `game_id` int(11) NOT NULL COMMENT 'game id',
  `game_contract_id` varchar(32) DEFAULT NULL COMMENT 'chaincode id',
  `title` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL,
  `gaming_seq` varchar(32) DEFAULT NULL COMMENT 'sequence of game　in progress',
  `status` int(4) DEFAULT '0',
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `remarks` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `gaming_seq` (`gaming_seq`)
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for public_random
-- ----------------------------
DROP TABLE IF EXISTS `public_random`;
CREATE TABLE `public_random` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'public_random id',
  `gaming_seq` varchar(32) DEFAULT NULL COMMENT 'sequence of game　in progress',
  `gaming_round` varchar(32) DEFAULT NULL COMMENT 'game round number',
  `public_input` varchar(255) NOT NULL,
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_gaming` (`gaming_seq`, `gaming_round`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for texas_poker
-- ----------------------------
DROP TABLE IF EXISTS `texas_poker`;
CREATE TABLE `texas_poker` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'texas poker id',
  `contract_id` varchar(32) DEFAULT NULL COMMENT 'chaincode id',
  `title` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL,
  `createtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `contract_id` (`contract_id`)
) ENGINE=InnoDB AUTO_INCREMENT=264 DEFAULT CHARSET=utf8;
