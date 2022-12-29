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

 Date: 10/10/2022 14:59:40
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ceping_sdk_audit_item
-- ----------------------------
DROP TABLE IF EXISTS `ceping_sdk_audit_item`;
CREATE TABLE `ceping_sdk_audit_item`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NULL DEFAULT NULL,
  `item_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `level` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `sort` int(11) NULL DEFAULT NULL,
  `status` int(11) NULL DEFAULT NULL,
  `score` int(11) NULL DEFAULT NULL,
  `solution` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `admin_setting` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `enable_item`(`item_key`, `status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 90 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ceping_sdk_audit_item
-- ----------------------------
INSERT INTO `ceping_sdk_audit_item` VALUES (1, 1, 'sdk_sec_infos', '基本信息', 'I', 1, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (2, 1, 'sdk_permission', '权限信息', 'I', 2, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (3, 1, 'sdk_behavior', '行为信息', 'I', 3, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (4, 1, 'sdk_tpSdk', '第三方SDK检测', 'I', 4, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (5, 1, 'sdk_words', '敏感词信息', 'I', 5, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (6, 1, 'sdk_scan_virus', 'SDK病毒扫描', 'H', 6, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (7, 1, 'sdk_res_apk', '资源文件中的Apk文件', 'L', 7, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (8, 1, 'sdk_custom_perms', '未保护的自定义权限风险', 'L', 8, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (9, 1, 'sdk_excessive_perm_announce', '权限过度声明', 'L', 9, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (10, 1, 'sdk_user_privacy_info', '用户隐私信息', 'M', 10, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (11, 2, 'sdk_soProtect', 'So文件破解风险', 'L', 1, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (12, 2, 'sdk_decompile', 'Java代码反编译风险', 'H', 2, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (13, 2, 'sdk_code_proguard', '代码未混淆风险', 'M', 3, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (14, 2, 'sdk_res_protect', '资源文件泄漏风险', 'L', 4, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (15, 2, 'sdk_only_use_java_code', '仅使用Java代码风险', 'L', 5, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (16, 2, 'sdk_shield', '加固壳识别', 'H', 6, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (17, 2, 'sdk_sign_cert', '应用签名算法不安全风险', 'L', 7, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (18, 2, 'sdk_testprop', '单元测试配置风险', 'M', 8, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (19, 3, 'sdk_url', '代码残留URL信息检测', 'L', 1, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (20, 3, 'sdk_sharedprefs', 'SharedPreferences数据全局可读写漏洞', 'M', 2, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (21, 3, 'sdk_getDir', 'getDir数据全局可读写漏洞', 'M', 3, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (22, 3, 'sdk_risk_method', '敏感函数调用风险', 'L', 4, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (23, 3, 'sdk_random', '随机数不安全使用漏洞', 'L', 5, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (24, 3, 'sdk_openFile', '全局可读写的内部文件漏洞', 'M', 6, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (25, 3, 'sdk_wv_savepwd', 'WebView明文存储密码风险', 'H', 7, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (26, 3, 'sdk_webview_file', 'Webview File域同源策略绕过', 'H', 8, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (27, 3, 'sdk_plaintext_cert', '明文数字证书风险', 'M', 9, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (28, 3, 'sdk_logapi', '调试日志函数调用风险', 'L', 10, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (29, 3, 'sdk_encrypt_risk', 'AES/DES加密算法不安全使用漏洞', 'L', 11, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (30, 3, 'sdk_rsa_risk', 'RSA加密算法不安全使用漏洞', 'M', 12, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (31, 3, 'sdk_key_risk', '密钥硬编码漏洞', 'H', 13, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (32, 3, 'sdk_webview_remote_debug', 'WebView远程调试风险', 'M', 14, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (33, 3, 'sdk_backup', '应用数据任意备份风险', 'M', 15, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (34, 3, 'sdk_dbglobalrw', '数据库文件任意读写漏洞', 'M', 16, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (35, 3, 'sdk_internal_storage_mode', 'Internal Storage数据全局可读写漏洞', 'M', 17, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (36, 3, 'sdk_sharedprefs_shareuserid', 'SharedUserId属性设置漏洞', 'M', 18, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (37, 3, 'sdk_ffmpeg_risk', 'ffmpge文件读取漏洞', 'H', 19, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (38, 3, 'sdk_java_debug_risk', 'Java层动态调试风险', 'M', 20, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (39, 3, 'sdk_clipboard', '剪切板敏感信息泄露漏洞', 'M', 21, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (40, 3, 'sdk_residua', '内网测试信息残留漏洞', 'L', 22, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (41, 3, 'sdk_sensitive_account_password', '残留账户密码信息检测', 'L', 23, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (42, 3, 'sdk_sensitive_phone', '残留手机号信息检测', 'L', 24, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (43, 3, 'sdk_residual_email', '残留Email信息', 'L', 25, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (44, 3, 'sdk_plaintext_secret_leak', '明文泄漏风险', 'M', 26, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (45, 3, 'sdk_cve_strand_hogg', 'StrandHogg漏洞', 'L', 27, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (46, 4, 'sdk_trans', 'HTTP传输数据风险', 'L', 1, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (47, 4, 'sdk_wv_sslerror', 'WebView绕过证书校验风险', 'M', 2, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (48, 4, 'sdk_x509trust', 'HTTPS未校验服务器证书漏洞', 'M', 3, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (49, 4, 'sdk_host_name', 'HTTPS未检验主机名漏洞', 'M', 4, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (50, 4, 'sdk_intermediator_risk', 'HTTPS允许任意主机名漏洞', 'M', 5, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (51, 4, 'sdk_network_env_check', '联网环境检测', 'I', 6, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (52, 4, 'sdk_oversea_server', '访问境外服务器检测', 'H', 7, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (53, 4, 'sdk_vpn_service', '启用vpn服务检测', 'I', 8, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (54, 5, 'sdk_screen_shots', '截屏攻击风险', 'M', 1, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (55, 5, 'sdk_kb_input', '输入监听风险', 'M', 2, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (56, 5, 'sdk_debug_certificate', '使用调试证书发布应用风险', 'L', 3, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (57, 6, 'sdk_receiver', '动态注册Receiver风险', 'M', 1, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (58, 6, 'sdk_pendingIntent', 'PendingIntent错误使用Intent风险', 'L', 2, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (59, 6, 'sdk_intent_hijack', 'Intent组件隐式调用风险', 'L', 3, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (60, 6, 'sdk_coms_activity', 'Activity组件导出风险', 'M', 4, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (61, 6, 'sdk_coms_service', 'Service组件导出风险', 'M', 5, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (62, 6, 'sdk_coms_receiver', 'Receiver组件导出风险', 'M', 6, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (63, 6, 'sdk_coms_provider', 'Provider组件导出风险', 'M', 7, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (64, 6, 'sdk_reflect', '反射调用风险', 'L', 8, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (65, 6, 'sdk_cve_gif_drawable', 'Android-gif-Drawable远程代码执行漏洞', 'H', 9, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (66, 6, 'sdk_intent_scheme_url', 'Intent Scheme URL攻击漏洞', 'M', 10, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (67, 6, 'sdk_fragment', 'Fragment注入攻击漏洞', 'M', 11, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (68, 7, 'sdk_h5_storage', 'Web Storage数据泄露风险', 'L', 1, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (69, 7, 'sdk_h5_websql', 'WebSQL漏洞', 'H', 2, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (70, 7, 'sdk_h5_innerhtml', 'innerHTML的漏洞', 'H', 3, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (71, 8, 'sdk_external_load_so', '从sdcard加载so风险', 'M', 1, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (72, 8, 'sdk_external_load_dex', '从sdcard加载dex风险', 'M', 2, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (73, 8, 'sdk_stack_protect', '未使用编译器堆栈保护技术', 'L', 3, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (74, 8, 'sdk_random_space', '未使用地址空间随机化技术风险', 'L', 4, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (75, 8, 'sdk_wv_inject', 'Webview远程代码执行漏洞', 'H', 5, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (76, 8, 'sdk_webview_fileurl', '“应用克隆”漏洞攻击风险', 'H', 6, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (77, 8, 'sdk_webview_hide_interface', '未移除有风险的Webview系统隐藏接口漏洞', 'H', 7, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (78, 8, 'sdk_unzip', 'zip文件解压目录遍历漏洞', 'M', 8, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (79, 8, 'sdk_anydown', '下载任意APK漏洞', 'H', 9, 1, 5, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (80, 8, 'sdk_risk_webBrowser', '不安全的浏览器调用', 'M', 10, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (81, 8, 'sdk_parasitic_push', '“寄生推”云控风险', 'L', 11, 1, 2, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (82, 8, 'sdk_nolauncher_service_risk', '启动隐藏服务风险', 'L', 12, 1, 3, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (83, 8, 'sdk_signaturev2', 'Janus签名机制漏洞', 'H', 13, 1, 6, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (84, 8, 'sdk_run_other_program', '运行其他可执行程序漏洞', 'M', 14, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (85, 9, 'sdk_sharedprefs_commit', 'SharedPreferences使用commit提交数据检测', 'I', 1, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (86, 9, 'sdk_global_exception', '全局异常处理', 'L', 2, 1, 4, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (87, 9, 'sdk_ip_address', 'IP地址检测', 'I', 3, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (88, 9, 'sdk_state_crypt_algorithm', '国密算法', 'I', 4, 1, 0, '', NULL);
INSERT INTO `ceping_sdk_audit_item` VALUES (89, 9, 'sdk_start_behavior', '自启行为', 'L', 5, 1, 4, '', NULL);

SET FOREIGN_KEY_CHECKS = 1;
