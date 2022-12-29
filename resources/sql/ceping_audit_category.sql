/*
 Navicat Premium Data Transfer

 Source Server         : 172.16.102.58
 Source Server Type    : MySQL
 Source Server Version : 80018
 Source Host           : 172.16.102.58:33060
 Source Schema         : ssp

 Target Server Type    : MySQL
 Target Server Version : 80018
 File Encoding         : 65001

 Date: 10/10/2022 16:15:47
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ceping_audit_category
-- ----------------------------
DROP TABLE IF EXISTS `ceping_audit_category`;
CREATE TABLE `ceping_audit_category`  (
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `id` int(11) NOT NULL,
  `category_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `category_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `category_sort` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`type`, `id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ceping_audit_category
-- ----------------------------
INSERT INTO `ceping_audit_category` VALUES ('ad', 1, 'sec', '自身安全', 1);
INSERT INTO `ceping_audit_category` VALUES ('ad', 2, 'src_check', '程序源文件安全', 2);
INSERT INTO `ceping_audit_category` VALUES ('ad', 3, 'local_data', '本地数据存储安全', 3);
INSERT INTO `ceping_audit_category` VALUES ('ad', 4, 'trans_safe', '通信数据传输安全', 4);
INSERT INTO `ceping_audit_category` VALUES ('ad', 5, 'identity_check', '身份认证安全', 5);
INSERT INTO `ceping_audit_category` VALUES ('ad', 6, 'internal_data_exchange', '内部数据交互安全', 6);
INSERT INTO `ceping_audit_category` VALUES ('ad', 7, 'component_security', '组件安全', 7);
INSERT INTO `ceping_audit_category` VALUES ('ad', 8, 'prevention_attack', '恶意攻击防范能力', 8);
INSERT INTO `ceping_audit_category` VALUES ('ad', 9, 'h5_hybrid', 'HTML5安全', 9);
INSERT INTO `ceping_audit_category` VALUES ('ad', 10, 'other_sdk_detect', '第三方SDK检测', 10);
INSERT INTO `ceping_audit_category` VALUES ('ad', 11, 'content_security', '内容安全', 11);
INSERT INTO `ceping_audit_category` VALUES ('ad', 12, 'optimization_suggestion', '优化建议', 12);
INSERT INTO `ceping_audit_category` VALUES ('hm', 1, 'hm_sec', '自身安全', 1);
INSERT INTO `ceping_audit_category` VALUES ('hm', 2, 'hm_src_check', '程序源文件安全', 2);
INSERT INTO `ceping_audit_category` VALUES ('hm', 3, 'hm_local_data', '本地数据存储安全', 3);
INSERT INTO `ceping_audit_category` VALUES ('hm', 4, 'hm_trans_safe', '通信数据传输安全', 4);
INSERT INTO `ceping_audit_category` VALUES ('hm', 5, 'hm_internal_data_exchange', '内部数据交互安全', 5);
INSERT INTO `ceping_audit_category` VALUES ('hm', 6, 'hm_prevention_attack', '恶意攻击防范能力', 6);
INSERT INTO `ceping_audit_category` VALUES ('hm', 7, 'hm_optimization_suggestion', '优化建议', 7);
INSERT INTO `ceping_audit_category` VALUES ('ios', 1, 'ios_sec', '自身安全', 1);
INSERT INTO `ceping_audit_category` VALUES ('ios', 2, 'ios_src_check', '二进制代码保护', 2);
INSERT INTO `ceping_audit_category` VALUES ('ios', 3, 'ios_local_data', '客户端数据存储安全', 3);
INSERT INTO `ceping_audit_category` VALUES ('ios', 4, 'ios_trans_safe', '数据传输安全', 4);
INSERT INTO `ceping_audit_category` VALUES ('ios', 5, 'ios_encryption', '加密算法及密码安全', 5);
INSERT INTO `ceping_audit_category` VALUES ('ios', 6, 'ios_app_sec', 'iOS应用安全规范', 6);
INSERT INTO `ceping_audit_category` VALUES ('ios', 7, 'ios_source_safe', '程序源文件安全', 7);
INSERT INTO `ceping_audit_category` VALUES ('ios', 8, 'ios_h5_hybrid', 'HTML5安全', 8);
INSERT INTO `ceping_audit_category` VALUES ('ios', 9, 'ios_other_sdk_detect', '第三方SDK检测', 9);
INSERT INTO `ceping_audit_category` VALUES ('ios', 10, 'ios_content_security', '内容安全', 10);
INSERT INTO `ceping_audit_category` VALUES ('mp', 1, 'mp_sec', '自身安全', 1);
INSERT INTO `ceping_audit_category` VALUES ('mp', 2, 'mp_comm_xfe', '通信传输安全检测', 2);
INSERT INTO `ceping_audit_category` VALUES ('mp', 3, 'mp_data_reveal', '数据泄漏检测', 3);
INSERT INTO `ceping_audit_category` VALUES ('mp', 4, 'mp_com_risk', '组件漏洞检测', 4);
INSERT INTO `ceping_audit_category` VALUES ('mp', 5, 'mp_unsafe_setting', 'HTTP不安全配置检测', 5);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 1, 'sdk_sec', '自身安全', 1);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 2, 'sdk_src_check', '程序源文件安全', 2);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 3, 'sdk_local_data', '本地数据存储安全', 3);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 4, 'sdk_trans_safe', '通信数据传输安全', 4);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 5, 'sdk_identity_check', '身份认证安全', 5);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 6, 'sdk_internal_data_exchange', '内部数据交互安全', 6);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 7, 'sdk_h5_hybrid', 'HTML5安全', 7);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 8, 'sdk_prevention_attack', '恶意攻击防范能力', 8);
INSERT INTO `ceping_audit_category` VALUES ('sdk', 9, 'sdk_optimization_suggestion', '优化建议', 9);

SET FOREIGN_KEY_CHECKS = 1;
