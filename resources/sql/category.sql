/*
 Navicat Premium Data Transfer

 Source Server         : 172.16.102.66
 Source Server Type    : MySQL
 Source Server Version : 80018
 Source Host           : 172.16.102.66:33060
 Source Schema         : middlegroundabc

 Target Server Type    : MySQL
 Target Server Version : 80018
 File Encoding         : 65001

 Date: 10/10/2022 17:12:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for category
-- ----------------------------
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `category_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `ce_ping_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 21 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of category
-- ----------------------------
INSERT INTO `category` VALUES (1, '自身安全', 'ad');
INSERT INTO `category` VALUES (2, '程序源文件安全', 'ad');
INSERT INTO `category` VALUES (3, '本地数据存储安全', 'ad');
INSERT INTO `category` VALUES (4, '通信数据传输安全', 'ad');
INSERT INTO `category` VALUES (5, '身份认证安全', 'ad');
INSERT INTO `category` VALUES (6, '内部数据交互安全', 'ad');
INSERT INTO `category` VALUES (7, '恶意攻击防范能力', 'ad');
INSERT INTO `category` VALUES (8, 'HTML5安全', 'ad');
INSERT INTO `category` VALUES (9, '第三方SDK检测', 'ad');
INSERT INTO `category` VALUES (10, '内容安全', 'ad');
INSERT INTO `category` VALUES (11, '优化建议', 'ad');
INSERT INTO `category` VALUES (12, '自身安全', 'ios');
INSERT INTO `category` VALUES (13, '二进制代码保护', 'ios');
INSERT INTO `category` VALUES (14, '客户端数据存储安全', 'ios');
INSERT INTO `category` VALUES (15, '数据传输安全', 'ios');
INSERT INTO `category` VALUES (16, '加密算法及密码安全', 'ios');
INSERT INTO `category` VALUES (17, 'iOS应用安全规范', 'ios');
INSERT INTO `category` VALUES (18, '程序源文件安全', 'ios');
INSERT INTO `category` VALUES (19, 'HTML5安全', 'ios');
INSERT INTO `category` VALUES (20, '第三方SDK检测', 'ios');
INSERT INTO `category` VALUES (21, '内容安全', 'ios');
INSERT INTO `category` VALUES (22, '自身安全', 'sdk');
INSERT INTO `category` VALUES (23, '程序源文件安全', 'sdk');
INSERT INTO `category` VALUES (24, '本地数据存储安全', 'sdk');
INSERT INTO `category` VALUES (25, '内部数据交互安全', 'sdk');
INSERT INTO `category` VALUES (26, '恶意攻击防范能力', 'sdk');

SET FOREIGN_KEY_CHECKS = 1;
