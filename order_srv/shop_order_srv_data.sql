/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80013
 Source Host           : localhost:3306
 Source Schema         : shop_order_srv

 Target Server Type    : MySQL
 Target Server Version : 80013
 File Encoding         : 65001

 Date: 30/07/2023 22:39:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ordergoods
-- ----------------------------
DROP TABLE IF EXISTS `ordergoods`;
CREATE TABLE `ordergoods`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `order` int(11) NULL DEFAULT NULL,
  `goods` int(11) NULL DEFAULT NULL,
  `goods_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `goods_image` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `goods_price` float NULL DEFAULT NULL,
  `nums` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ordergoods_order`(`order`) USING BTREE,
  INDEX `idx_ordergoods_goods`(`goods`) USING BTREE,
  INDEX `idx_ordergoods_goods_name`(`goods_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ordergoods
-- ----------------------------
INSERT INTO `ordergoods` VALUES (3, '2023-07-19 22:43:09.251', '2023-07-19 22:43:09.251', NULL, 0, 11, 422, '西州蜜瓜25号哈密瓜 2粒装 单果1.25kg以上 新鲜水果', 'https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73', 36.9, 1);
INSERT INTO `ordergoods` VALUES (4, '2023-07-30 22:29:58.641', '2023-07-30 22:29:58.641', NULL, 0, 12, 422, '西州蜜瓜25号哈密瓜 2粒装 单果1.25kg以上 新鲜水果', 'https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/c3dee23a62efe14bbd4fc2c70046dc73', 36.9, 1);

-- ----------------------------
-- Table structure for orderinfo
-- ----------------------------
DROP TABLE IF EXISTS `orderinfo`;
CREATE TABLE `orderinfo`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL,
  `order_sn` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `pay_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'alipay(支付宝)， wechat(微信)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)',
  `trade_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '交易号',
  `order_mount` float NULL DEFAULT NULL,
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支持时间',
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收获地址',
  `signer_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货人姓名',
  `singer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '收货人联系方式',
  `post` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '留言信息',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_orderinfo_user`(`user`) USING BTREE,
  INDEX `idx_orderinfo_order_sn`(`order_sn`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orderinfo
-- ----------------------------
INSERT INTO `orderinfo` VALUES (11, '2023-07-19 22:43:08.737', '2023-07-19 22:43:08.737', NULL, 0, 0, '202371922437366481003189', '', '', '', 36.9, NULL, '浙江省', 'piwriw', '18888888', '请尽快发货');
INSERT INTO `orderinfo` VALUES (12, '2023-07-30 22:29:58.639', '2023-07-30 22:29:58.639', NULL, 0, 0, '202373022296385431003122', '', '', '', 36.9, NULL, '浙江省', 'test', '18268958295', '18268958295');

-- ----------------------------
-- Table structure for shoppingcart
-- ----------------------------
DROP TABLE IF EXISTS `shoppingcart`;
CREATE TABLE `shoppingcart`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL,
  `goods` int(11) NULL DEFAULT NULL,
  `nums` int(11) NULL DEFAULT NULL,
  `checked` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_shoppingcart_user`(`user`) USING BTREE,
  INDEX `idx_shoppingcart_goods`(`goods`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shoppingcart
-- ----------------------------
INSERT INTO `shoppingcart` VALUES (1, '2023-07-19 21:37:02.855', '2023-07-30 13:27:56.967', '2023-07-30 19:50:58.328', 0, 31, 422, 4, 1);
INSERT INTO `shoppingcart` VALUES (2, '2023-07-30 19:52:33.607', '2023-07-30 19:56:54.586', '2023-07-30 22:29:58.650', 0, 31, 422, 1, 1);

SET FOREIGN_KEY_CHECKS = 1;
