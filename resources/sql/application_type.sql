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

 Date: 18/11/2021 15:39:59
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for application_type
-- ----------------------------
DROP TABLE IF EXISTS `application_type`;
CREATE TABLE `application_type`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT,
  `app_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of application_type
-- ----------------------------
INSERT INTO `application_type` VALUES (1, 'android加固');
INSERT INTO `application_type` VALUES (2, 'h5加固');

SET FOREIGN_KEY_CHECKS = 1;
