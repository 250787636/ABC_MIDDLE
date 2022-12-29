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

 Date: 10/10/2022 16:16:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ceping_ios_audit_item
-- ----------------------------
DROP TABLE IF EXISTS `ceping_ios_audit_item`;
CREATE TABLE `ceping_ios_audit_item`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NULL DEFAULT NULL,
  `item_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `score` int(11) NULL DEFAULT NULL,
  `sort` int(11) NULL DEFAULT NULL,
  `status` tinyint(4) NULL DEFAULT NULL,
  `solution` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 46 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ceping_ios_audit_item
-- ----------------------------
INSERT INTO `ceping_ios_audit_item` VALUES (1, 1, 'ios_sec_infos', '基本信息', 'I', 0, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (2, 1, 'ios_sec_perms', '权限信息', 'I', 0, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (3, 1, 'ios_sec_behavior', '行为信息', 'I', 0, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (4, 1, 'ios_check_certificate_type', '证书类型检测', 'L', 2, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (5, 1, 'ios_appstore_risk', '无法上架Appstore风险', 'L', 2, 5, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (6, 1, 'ios_compiler_architecture', '编译架构检测', 'I', 0, 6, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (7, 2, 'ios_aud_code_proguard', '代码未混淆风险', 'M', 4, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (8, 2, 'ios_compile_pie', '未使用地址空间随机化技术风险', 'L', 3, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (9, 2, 'ios_compile_ssp', '未使用编译器堆栈保护技术风险', 'L', 3, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (10, 2, 'ios_third_library_inject', '注入攻击风险', 'H', 6, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (11, 2, 'ios_maco_format', '可执行文件被篡改风险', 'M', 6, 5, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (12, 3, 'ios_aud_debug', '动态调试攻击风险', 'H', 6, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (13, 3, 'ios_cvs_keyboard_hijack', '输入监听风险', 'M', 4, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (14, 3, 'ios_cvs_debug_info', '调试日志函数调用风险', 'L', 2, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (15, 3, 'ios_cvs_webView_access_file', 'Webview组件跨域访问风险', 'H', 6, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (16, 3, 'ios_cvs_prison_break', '越狱设备运行风险', 'L', 3, 5, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (17, 3, 'ios_cvs_sqlite_risk', '数据库明文存储风险', 'M', 4, 6, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (18, 3, 'ios_cvs_profile_leakage', '配置文件信息明文存储风险', 'H', 5, 7, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (19, 3, 'ios_sec_other_frameworks', '动态库信息泄露风险', 'L', 3, 8, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (20, 4, 'ios_cvs_http_protocol', 'HTTP传输数据风险', 'L', 3, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (21, 4, 'ios_cvs_https_auth', 'HTTPS未校验服务器证书漏洞', 'L', 3, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (22, 4, 'ios_cvs_url_schemes', 'URL Schemes劫持漏洞', 'M', 4, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (23, 4, 'ios_network_env_check', '联网环境检测', 'I', 0, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (24, 5, 'ios_weak_encryption', 'AES/DES加密算法不安全使用漏洞', 'L', 3, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (25, 5, 'ios_cvs_weak_hash', '弱哈希算法使用漏洞', 'L', 3, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (26, 5, 'ios_cvs_random_risk', '随机数不安全使用漏洞', 'L', 3, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (27, 6, 'ios_cvs_xcode_ghost', 'XcodeGhost感染漏洞', 'H', 6, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (28, 6, 'ios_cvs_high_risk_api', '不安全的API函数引用风险', 'H', 5, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (29, 6, 'ios_cvs_private_api', 'Private Methods使用检测', 'H', 5, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (30, 6, 'ios_cvs_zip_down', 'ZipperDown解压漏洞', 'M', 4, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (31, 6, 'ios_cvs_iback_door', 'iBackDoor控制漏洞', 'H', 6, 5, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (32, 6, 'ios_compile_arc', '未使用自动管理内存技术风险', 'L', 3, 6, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (33, 6, 'ios_compile_arc_api', '内存分配函数不安全风险', 'L', 3, 7, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (34, 6, 'ios_custom_method_long', '自定义函数逻辑过于复杂风险', 'M', 6, 8, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (35, 7, 'ios_str_leakage', '明文字符串泄露风险', 'L', 2, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (36, 7, 'ios_explicit_syscall', '外部函数显式调用风险', 'L', 2, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (37, 7, 'ios_syscall', '系统调用暴露风险', 'L', 2, 3, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (38, 7, 'ios_create_exec_mem', '创建可执行权限内存风险', 'M', 4, 4, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (39, 7, 'ios_resign', '篡改和二次打包风险', 'H', 6, 5, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (40, 7, 'ios_sql_exec_code', 'SQLite内存破坏漏洞', 'M', 4, 6, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (41, 7, 'ios_str_format', '格式化字符串漏洞', 'M', 4, 7, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (42, 8, 'h5_ios_storage', 'Web Storage数据泄露风险', 'L', 0, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (43, 8, 'h5_ios_innerhtml', 'InnerHTML的XSS攻击漏洞', 'H', 5, 2, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (44, 9, 'ios_sec_other_sdk', '第三方SDK检测', 'I', 0, 1, 1, '');
INSERT INTO `ceping_ios_audit_item` VALUES (45, 10, 'ios_cvs_words', '敏感词', 'I', 0, 1, 1, '');

SET FOREIGN_KEY_CHECKS = 1;
