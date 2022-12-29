/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : localhost:3306
 Source Schema         : middlegroundabc

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 09/12/2021 15:03:33
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `user_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `department_id` bigint(0) UNSIGNED NULL DEFAULT NULL,
  `account_level` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `job_title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `last_login_time` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `client_token` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `token` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `is_admin` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, '2021-12-09 15:02:00.903', '2021-12-09 15:02:00.903', NULL, 'root', 0, '超级管理员', '超级管理员', '2021-12-09 15:02:00','','0caa659aef0c17bc0bd89efc609964f6', 1);

SET FOREIGN_KEY_CHECKS = 1;
