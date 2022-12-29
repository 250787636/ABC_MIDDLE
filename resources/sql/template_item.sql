/*
 Navicat Premium Data Transfer

 Source Server         : 172.16.38.215
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : 172.16.38.215:33306
 Source Schema         : middlegroundabc

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 16/12/2021 11:27:49
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for template_item
-- ----------------------------
DROP TABLE IF EXISTS `template_item`;
CREATE TABLE `template_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `category_id` bigint(20) DEFAULT NULL,
  `item_key` longtext,
  `audit_name` longtext,
  `is_dynamic` tinyint(1) DEFAULT NULL,
  `status` bigint(20) DEFAULT NULL,
  `ce_ping_type` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=182 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of template_item
-- ----------------------------
BEGIN;
INSERT INTO `template_item` VALUES (1, 8, 'h5_ios_storage', 'Web Storage数据泄露风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (2, 8, 'h5_ios_innerhtml', 'InnerHTML的XSS攻击漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (3, 6, 'ios_cvs_xcode_ghost', 'XcodeGhost感染漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (4, 6, 'ios_cvs_high_risk_api', '不安全的API函数引用风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (5, 6, 'ios_cvs_private_api', 'Private Methods使用检测', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (6, 6, 'ios_cvs_zip_down', 'ZipperDown解压漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (7, 6, 'ios_cvs_iback_door', 'iBackDoor控制漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (8, 6, 'ios_compile_arc', '未使用自动管理内存技术风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (9, 6, 'ios_compile_arc_api', '内存分配函数不安全风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (10, 6, 'ios_custom_method_long', '自定义函数逻辑过于复杂风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (11, 2, 'ios_third_library_inject', '注入攻击风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (12, 2, 'ios_maco_format', '可执行文件被篡改风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (13, 2, 'ios_aud_code_proguard', '代码未混淆风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (14, 2, 'ios_compile_pie', '未使用地址空间随机化技术风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (15, 2, 'ios_compile_ssp', '未使用编译器堆栈保护技术风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (16, 5, 'ios_weak_encryption', 'AES/DES加密算法不安全使用漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (17, 5, 'ios_cvs_weak_hash', '弱哈希算法使用漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (18, 5, 'ios_cvs_random_risk', '随机数不安全使用漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (19, 3, 'ios_aud_debug', '动态调试攻击风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (20, 3, 'ios_cvs_keyboard_hijack', '输入监听风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (21, 3, 'ios_cvs_debug_info', '调试日志函数调用风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (22, 3, 'ios_cvs_webView_access_file', 'Webview组件跨域访问风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (23, 3, 'ios_cvs_prison_break', '越狱设备运行风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (24, 3, 'ios_cvs_sqlite_risk', '数据库明文存储风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (25, 3, 'ios_cvs_profile_leakage', '配置文件信息明文存储风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (26, 3, 'ios_sec_other_frameworks', '动态库信息泄露风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (27, 4, 'ios_cvs_http_protocol', 'HTTP传输数据风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (28, 4, 'ios_cvs_https_auth', 'HTTPS未校验服务器证书漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (29, 4, 'ios_cvs_url_schemes', 'URL Schemes劫持漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (31, 7, 'ios_resign', '篡改和二次打包风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (32, 7, 'ios_create_exec_mem', '创建可执行权限内存风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (33, 7, 'ios_sql_exec_code', 'SQLite内存破坏漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (34, 7, 'ios_str_format', '格式化字符串漏洞', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (35, 7, 'ios_str_leakage', '明文字符串泄露风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (36, 7, 'ios_explicit_syscall', '外部函数显式调用风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (37, 7, 'ios_syscall', '系统调用暴露风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (38, 1, 'ios_sec_perms', '权限信息', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (39, 1, 'ios_sec_behavior', '行为信息', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (40, 1, 'ios_check_certificate_type', '证书类型检测', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (41, 1, 'ios_appstore_risk', '无法上架Appstore风险', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (43, 15, 'aud_webview_fileurl', '“应用克隆”漏洞攻击风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (44, 15, 'aud_inject', '动态注入攻击风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (45, 15, 'cvs_wv_inject', 'Webview远程代码执行漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (46, 15, 'aud_webview_hide_interface', '未移除有风险的Webview系统隐藏接口漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (47, 15, 'aud_unzip', 'zip文件解压目录遍历漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (48, 15, 'cvs_anydown', '下载任意apk漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (49, 15, 'cvs_refuse_service', '拒绝服务攻击漏洞', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (50, 15, 'aud_sdcard_loaddex', '从sdcard加载dex风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (51, 15, 'aud_sdcard_loadso', '从sdcard加载so风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (52, 15, 'aud_stack_protect', '未使用编译器堆栈保护技术风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (53, 15, 'aud_random_space', '未使用地址空间随机化技术风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (55, 15, 'aud_root_device', 'Root设备运行风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (56, 15, 'aud_risk_webBrowser', '不安全的浏览器调用漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (58, 11, 'aud_savepwd', 'Webview明文存储密码风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (59, 11, 'aud_webview_file', 'Webview File同源策略绕过漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (60, 11, 'aud_cert', '明文数字证书风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (61, 11, 'aud_logapi', '调试日志函数调用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (62, 11, 'cvs_sqlinject', '数据库注入漏洞', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (63, 11, 'cvs_encrypt_risk', 'AES/DES加密方法不安全使用漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (64, 11, 'cvs_rsa_risk', 'RSA加密算法不安全使用漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (65, 11, 'aud_key_risk', '密钥硬编码漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (66, 11, 'aud_attack', '动态调试攻击风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (67, 11, 'aud_webview_remote_debug', 'Webview远程调试风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (68, 11, 'aud_backup', '应用数据任意备份风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (69, 11, 'aud_sensapi', '敏感函数调用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (70, 11, 'aud_dbglobalrw', '数据库文件任意读写漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (71, 11, 'cvs_globalrw', '全局可读写的内部文件漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (72, 11, 'cvs_sharedprefs', 'SharedPreferences数据全局可读写漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (73, 11, 'cvs_sharedprefs_shareuserid', 'SharedUserId属性设置漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (74, 11, 'cvs_internal_storage_mode', 'Internal Storage数据全局可读写漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (75, 11, 'aud_get_dir', 'getDir数据全局可读写漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (76, 11, 'aud_ffmpeg_risk', 'FFmpeg文件读取漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (77, 11, 'aud_debug', 'Java层动态调试风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (78, 11, 'aud_clipboard', '剪切板敏感信息泄露漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (79, 11, 'cvs_residua', '内网测试信息残留漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (80, 11, 'cvs_random', '随机数不安全使用漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (81, 11, 'aud_url', '代码残留URL信息检测', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (82, 11, 'aud_sensitive_account_password', '残留账户密码信息检测', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (83, 11, 'aud_sensitive_phone', '残留手机号信息检测', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (87, 7, 'aud_shield', '加固壳识别', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (88, 7, 'aud_decompile', 'Java代码反编译风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (89, 7, 'aud_so_protect', 'So文件破解风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (90, 7, 'aud_tamper', '篡改和二次打包风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (91, 7, 'aud_signaturev2', 'Janus签名机制漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (92, 7, 'aud_res_protect', '资源文件泄露风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (93, 7, 'aud_signature', '应用签名未校验风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (94, 7, 'aud_code_proguard', '代码未混淆风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (95, 7, 'aud_certificate', '使用调试证书发布应用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (96, 7, 'aud_JniRegisterNatives_risk', '仅使用Java代码风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (97, 7, 'aud_nolauncher_service_risk', '启动隐藏服务风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (98, 7, 'aud_sign_cert', '应用签名算法不安全风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (99, 7, 'aud_testprop', '单元测试配置风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (100, 10, 'aud_shield', '加固壳识别', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (101, 10, 'aud_decompile', 'Java代码反编译风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (102, 10, 'aud_so_protect', 'So文件破解风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (103, 10, 'aud_tamper', '篡改和二次打包风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (104, 10, 'aud_signaturev2', 'Janus签名机制漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (105, 10, 'aud_res_protect', '资源文件泄露风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (106, 10, 'aud_signature', '应用签名未校验风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (107, 10, 'aud_code_proguard', '代码未混淆风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (108, 10, 'aud_certificate', '使用调试证书发布应用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (109, 10, 'aud_JniRegisterNatives_risk', '仅使用Java代码风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (110, 10, 'aud_nolauncher_service_risk', '启动隐藏服务风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (111, 10, 'aud_sign_cert', '应用签名算法不安全风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (113, 1, 'sec_perms', '权限信息', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (114, 1, 'sec_behavior', '行为信息', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (115, 1, 'sec_virus', '病毒扫描', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (116, 1, 'aud_res_apk', '资源文件中的Apk文件', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (117, 1, 'sec_excessive_perm_announce', '权限过度声明风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (118, 1, 'sec_custom_perms', '未保护的自定义权限风险检测', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (119, 9, 'sec_perms', '权限信息', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (120, 9, 'sec_behavior', '行为信息', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (121, 9, 'sec_virus', '病毒扫描', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (122, 9, 'aud_res_apk', '资源文件中的Apk文件', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (123, 9, 'sec_other_sdk', '第三方SDK检测', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (124, 9, 'sec_words', '敏感词信息', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (125, 13, 'aud_uihijack', '界面劫持风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (126, 13, 'aud_kb_input', '输入监听风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (127, 13, 'aud_screen_shots', '截屏攻击风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (128, 12, 'aud_trans', 'HTTP传输数据风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (129, 12, 'cvs_x509trust', 'HTTPS未校验服务器证书漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (130, 12, 'aud_host_name', 'HTTPS未校验主机名漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (131, 12, 'aud_intermediator_risk', 'HTTPS允许任意主机名漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (132, 12, 'cvs_wv_sslerror', 'Webview绕过证书校验漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (133, 12, 'cvs_packet_capture', 'HTTP报文信息泄漏风险', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (137, 8, 'h5_storage', 'Web Storage数据泄露风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (138, 8, 'h5_websql', 'WebSQL注入漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (139, 8, 'h5_innerhtml', 'InnerHTML的XSS攻击漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (140, 16, 'h5_storage', 'Web Storage数据泄露风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (141, 16, 'h5_websql', 'WebSQL注入漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (142, 16, 'h5_innerhtml', 'InnerHTML的XSS攻击漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (143, 14, 'aud_register_receiver', '动态注册Receiver风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (144, 14, 'cvs_dataleak', 'Content Provider数据泄露漏洞', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (145, 14, 'aud_coms_activity', 'Activity组件导出风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (146, 14, 'aud_coms_service', 'Service组件导出风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (147, 14, 'aud_coms_receiver', 'Broadcast Receiver组件导出风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (148, 14, 'aud_coms_provider', 'Content Provider组件导出风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (149, 14, 'cvs_local_port', '本地端口开放越权漏洞', 1, 0, 'ad');
INSERT INTO `template_item` VALUES (150, 14, 'aud_pendingintent', 'PendingIntent错误使用Intent风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (151, 14, 'aud_component_hijack', 'Intent组件隐式调用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (152, 14, 'cvs_intent_risk', 'Intent Scheme URL攻击漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (153, 14, 'cvs_fragment_risk', 'Fragment注入攻击漏洞', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (154, 14, 'aud_reflect', '反射调用风险', 0, 0, 'ad');
INSERT INTO `template_item` VALUES (156, 21, 'mp_cors_risk', '不安全的跨域资源共享配置风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (157, 21, 'mp_hsts_risk', 'HTTP严格传输安全策略检测', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (158, 20, 'mp_shell_shock', '破壳（ShellShock）漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (159, 20, 'mp_heart_bleed', '心脏滴血漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (160, 20, 'mp_ssl_poodle', 'SSL POODLE漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (161, 20, 'mp_ssl_drown', 'SSL DROWN攻击漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (162, 20, 'mp_ssl_freak', 'SSL FREAK漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (163, 1, 'mp_sec_perms', '权限信息', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (164, 1, 'mp_open_port', '服务器开放端口检测', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (165, 9, 'mp_sec_perms', '权限信息', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (166, 9, 'mp_open_port', '服务器开放端口检测', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (167, 17, 'mp_sec_perms', '权限信息', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (168, 17, 'mp_open_port', '服务器开放端口检测', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (169, 18, 'mp_ssl_host_error', 'SSL证书主机名错误漏洞', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (170, 18, 'mp_ssl_self_sign', '自签名SSL证书使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (171, 18, 'mp_ssl_no_trust', '使用不被信任的SSL证书风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (172, 18, 'mp_unuse_https', 'HTTP传输数据风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (173, 18, 'mp_ssl_cs', 'SSL使用中等强度密码套件风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (174, 18, 'mp_ssl_rc4', 'SSL RC4密码套件使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (175, 18, 'mp_http_put', 'HTTP PUT方法使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (176, 18, 'mp_http_trace', 'HTTP TRACE方法使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (177, 18, 'mp_http_move', 'HTTP MOVE方法使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (178, 18, 'mp_http_copy', 'HTTP COPY方法使用风险', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (179, 18, 'mp_http_request', 'HTTP网络请求列表', 0, 0, 'mp');
INSERT INTO `template_item` VALUES (180, 1, 'ios_cvs_words', '敏感词信息', 0, 0, 'ios');
INSERT INTO `template_item` VALUES (181, 1, 'ios_sec_other_sdk', '第三方SDK检测', 0, 0, 'ios');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
