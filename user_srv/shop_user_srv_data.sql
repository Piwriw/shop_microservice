/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80013
 Source Host           : localhost:3306
 Source Schema         : shop_user_srv

 Target Server Type    : MySQL
 Target Server Version : 80013
 File Encoding         : 65001

 Date: 26/06/2023 17:42:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `delete_at` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nick_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `birthday` datetime NULL DEFAULT NULL,
  `gender` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'male' COMMENT 'male男，female女',
  `role` int(11) NULL DEFAULT 1 COMMENT '1表示普通成员，2表示管理员',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 51 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (31, '2023-06-13 17:26:08.099', '2023-06-13 17:26:08.099', NULL, 0, '18268958295', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw0', '2023-06-13 17:26:08', 'male', 2);
INSERT INTO `users` VALUES (32, '2023-06-13 17:26:08.111', '2023-06-13 17:26:08.111', NULL, 0, '1111111', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw1', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (33, '2023-06-13 17:26:08.113', '2023-06-13 17:26:08.113', NULL, 0, '1111112', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw2', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (34, '2023-06-13 17:26:08.116', '2023-06-13 17:26:08.116', NULL, 0, '1111113', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw3', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (35, '2023-06-13 17:26:08.117', '2023-06-13 17:26:08.117', NULL, 0, '1111114', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw4', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (36, '2023-06-13 17:26:08.119', '2023-06-13 17:26:08.119', NULL, 0, '1111115', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw5', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (37, '2023-06-13 17:26:08.121', '2023-06-13 17:26:08.121', NULL, 0, '1111116', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw6', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (38, '2023-06-13 17:26:08.123', '2023-06-13 17:26:08.123', NULL, 0, '1111117', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw7', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (39, '2023-06-13 17:26:08.125', '2023-06-13 17:26:08.125', NULL, 0, '1111118', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw8', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (40, '2023-06-13 17:26:08.128', '2023-06-13 17:26:08.128', NULL, 0, '1111119', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw9', '2023-06-13 17:26:08', 'male', 1);
INSERT INTO `users` VALUES (41, '2023-06-16 09:21:41.205', '2023-06-16 09:21:41.205', NULL, 0, '1111110', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw0', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (42, '2023-06-16 09:21:41.216', '2023-06-16 09:21:41.216', NULL, 0, '1111111', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw1', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (43, '2023-06-16 09:21:41.218', '2023-06-16 09:21:41.218', NULL, 0, '1111112', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw2', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (44, '2023-06-16 09:21:41.219', '2023-06-16 09:21:41.219', NULL, 0, '1111113', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw3', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (45, '2023-06-16 09:21:41.221', '2023-06-16 09:21:41.221', NULL, 0, '1111114', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw4', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (46, '2023-06-16 09:21:41.223', '2023-06-16 09:21:41.223', NULL, 0, '1111115', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw5', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (47, '2023-06-16 09:21:41.225', '2023-06-16 09:21:41.225', NULL, 0, '1111116', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw6', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (48, '2023-06-16 09:21:41.227', '2023-06-16 09:21:41.227', NULL, 0, '1111117', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw7', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (49, '2023-06-16 09:21:41.229', '2023-06-16 09:21:41.229', NULL, 0, '1111118', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw8', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (50, '2023-06-16 09:21:41.230', '2023-06-16 09:21:41.230', NULL, 0, '1111119', 'e36b053c09461511bf5321a72afd9ccf', 'piwriw9', '2023-06-16 09:21:41', 'male', 1);
INSERT INTO `users` VALUES (51, '2023-06-17 16:25:02.701', '2023-06-17 16:25:02.701', NULL, 0, '18268958265', '4a0be41b614932ff93b7c245d5de551d', 'huan', NULL, 'male', 1);

SET FOREIGN_KEY_CHECKS = 1;
