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

 Date: 10/10/2022 16:16:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ceping_ad_audit_item
-- ----------------------------
DROP TABLE IF EXISTS `ceping_ad_audit_item`;
CREATE TABLE `ceping_ad_audit_item`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NULL DEFAULT NULL,
  `item_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `score` int(11) NULL DEFAULT NULL,
  `is_dynamic` tinyint(4) NULL DEFAULT NULL,
  `sort` int(11) NULL DEFAULT NULL,
  `status` bigint(20) NULL DEFAULT NULL,
  `solution` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `enable_item`(`item_key`, `status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 131 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ceping_ad_audit_item
-- ----------------------------
INSERT INTO `ceping_ad_audit_item` VALUES (1, 1, 'sec_infos', '基本信息', 'I', 0, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (2, 1, 'sec_perms', '权限信息', 'I', 0, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (3, 1, 'sec_behavior', '行为信息', 'I', 0, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (4, 1, 'sec_virus', '病毒扫描', 'H', 6, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (5, 1, 'aud_res_apk', '资源文件中的Apk文件', 'L', 0, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (6, 1, 'sec_excessive_perm_announce', '权限过度声明风险', 'L', 2, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (7, 1, 'sec_custom_perms', '未保护的自定义权限风险检测', 'L', 2, 0, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (8, 1, 'aud_app_testonly', '应用测试模式发布风险', 'L', 0, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (9, 1, 'aud_cert_markets', '来源安全检测', 'I', 0, 0, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (10, 1, 'sec_dangerous_perms', '应用权限安全', 'I', 0, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (11, 1, 'sec_control_perms', '控制力安全', 'I', 0, 0, 11, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (12, 2, 'aud_shield', '加固壳识别', 'H', 5, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (13, 2, 'aud_decompile', 'Java代码反编译风险', 'H', 6, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (14, 2, 'aud_so_protect', 'So文件破解风险', 'L', 3, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (15, 2, 'aud_tamper', '篡改和二次打包风险', 'H', 6, 1, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (16, 2, 'aud_signaturev2', 'Janus签名机制漏洞', 'H', 6, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (17, 2, 'aud_res_protect', '资源文件泄露风险', 'L', 2, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (18, 2, 'aud_signature', '应用签名未校验风险', 'H', 10, 1, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (19, 2, 'aud_code_proguard', '代码未混淆风险', 'M', 4, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (20, 2, 'aud_certificate', '使用调试证书发布应用风险', 'L', 3, 0, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (21, 2, 'aud_JniRegisterNatives_risk', '仅使用Java代码风险', 'L', 3, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (22, 2, 'aud_nolauncher_service_risk', '启动隐藏服务风险', 'L', 3, 0, 11, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (23, 2, 'aud_sign_cert', '应用签名算法不安全风险', 'L', 3, 0, 12, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (24, 2, 'aud_testprop', '单元测试配置风险', 'M', 4, 0, 13, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (25, 2, 'aud_res_xml_protect', 'xml资源文件泄露风险', 'L', 0, 0, 14, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (26, 2, 'cvs_umeng_sdk', '友盟SDK越权漏洞', 'M', 0, 0, 15, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (27, 2, 'aud_malicious_urls', '恶意URL检测', 'L', 0, 0, 16, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (28, 3, 'aud_savepwd', 'Webview明文存储密码风险', 'H', 5, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (29, 3, 'aud_webview_file', 'Webview File同源策略绕过漏洞', 'H', 5, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (30, 3, 'aud_cert', '明文数字证书风险', 'M', 4, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (31, 3, 'aud_logapi', '调试日志函数调用风险', 'L', 2, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (32, 3, 'cvs_sqlinject', '数据库注入漏洞', 'H', 5, 1, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (33, 3, 'cvs_encrypt_risk', 'AES/DES加密方法不安全使用漏洞', 'L', 3, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (34, 3, 'cvs_rsa_risk', 'RSA加密算法不安全使用漏洞', 'M', 4, 0, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (35, 3, 'aud_key_risk', '密钥硬编码漏洞', 'H', 5, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (36, 3, 'aud_attack', '动态调试攻击风险', 'H', 6, 1, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (37, 3, 'aud_webview_remote_debug', 'Webview远程调试风险', 'M', 4, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (38, 3, 'aud_backup', '应用数据任意备份风险', 'M', 4, 0, 11, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (39, 3, 'aud_sensapi', '敏感函数调用风险', 'L', 2, 0, 12, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (40, 3, 'aud_dbglobalrw', '数据库文件任意读写漏洞', 'M', 4, 0, 13, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (41, 3, 'cvs_globalrw', '全局可读写的内部文件漏洞', 'M', 4, 0, 14, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (42, 3, 'cvs_sharedprefs', 'SharedPreferences数据全局可读写漏洞', 'M', 4, 0, 15, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (43, 3, 'cvs_sharedprefs_shareuserid', 'SharedUserId属性设置漏洞', 'M', 4, 0, 16, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (44, 3, 'cvs_internal_storage_mode', 'Internal Storage数据全局可读写漏洞', 'M', 4, 0, 17, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (45, 3, 'aud_get_dir', 'getDir数据全局可读写漏洞', 'M', 4, 0, 18, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (46, 3, 'aud_ffmpeg_risk', 'FFmpeg文件读取漏洞', 'H', 5, 0, 19, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (47, 3, 'aud_debug', 'Java层动态调试风险', 'M', 4, 0, 20, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (48, 3, 'aud_clipboard', '剪切板敏感信息泄露漏洞', 'M', 4, 0, 21, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (49, 3, 'cvs_residua', '内网测试信息残留漏洞', 'L', 2, 0, 22, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (50, 3, 'cvs_random', '随机数不安全使用漏洞', 'L', 3, 0, 23, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (51, 3, 'aud_url', '代码残留URL信息检测', 'L', 2, 0, 24, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (52, 3, 'aud_sensitive_account_password', '残留账户密码信息检测', 'L', 2, 0, 25, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (53, 3, 'aud_sensitive_phone', '残留手机号信息检测', 'L', 2, 0, 26, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (54, 3, 'aud_residual_email', '残留Email信息检测', 'L', 2, 0, 27, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (55, 3, 'sec_plaintext_secret_leak', '明文泄漏风险', 'M', 6, 0, 28, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (56, 3, 'cve_strand_hogg', 'StrandHogg漏洞', 'L', 2, 0, 29, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (57, 3, 'aud_storage_operation', '储存卡的操作行为', 'L', 0, 0, 30, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (58, 3, 'aud_private_ip_expose', '私有IP地址暴露', 'L', 0, 0, 31, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (59, 3, 'cvs_ecb_encrypt_risk', 'ECB模式的AES/DEA加密方法不安全使用漏洞', 'L', 0, 0, 32, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (60, 3, 'cvs_ofb_encrypt_risk', 'OFB模式的AES/DEA加密方法不安全使用漏洞', 'L', 0, 0, 33, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (61, 3, 'aud_str_convert_risk', '字节数组与字符串转换风险', 'L', 0, 0, 34, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (62, 3, 'aud_define_key_size_risk', '用户控制的Key长度', 'L', 0, 0, 35, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (63, 3, 'aud_unsafe_iv_risk', '不安全的初始化向量', 'L', 0, 0, 36, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (64, 4, 'aud_trans', 'HTTP传输数据风险', 'L', 2, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (65, 4, 'cvs_x509trust', 'HTTPS未校验服务器证书漏洞', 'M', 4, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (66, 4, 'aud_host_name', 'HTTPS未校验主机名漏洞', 'M', 4, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (67, 4, 'aud_intermediator_risk', 'HTTPS允许任意主机名漏洞', 'M', 4, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (68, 4, 'cvs_wv_sslerror', 'Webview绕过证书校验漏洞', 'M', 4, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (69, 4, 'cvs_packet_capture', 'HTTP报文信息泄漏风险', 'L', 2, 1, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (70, 4, 'aud_network_env_check', '联网环境检测', 'I', 0, 0, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (71, 4, 'cvs_vpn_service', '启用VPN服务检测', 'I', 0, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (72, 4, 'aud_oversea_server', '访问境外服务器风险', 'L', 2, 0, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (73, 4, 'aud_socket_security', '通信套接字安全', 'M', 0, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (74, 5, 'aud_uihijack', '界面劫持风险', 'M', 4, 1, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (75, 5, 'aud_kb_input', '输入监听风险', 'M', 4, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (76, 5, 'aud_screen_shots', '截屏攻击风险', 'M', 4, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (77, 6, 'aud_register_receiver', '动态注册Receiver风险', 'M', 4, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (78, 6, 'cvs_dataleak', 'Content Provider数据泄露漏洞', 'H', 5, 1, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (79, 6, 'cvs_local_port', '本地端口开放越权漏洞', 'M', 4, 1, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (80, 6, 'aud_pendingintent', 'PendingIntent错误使用Intent风险', 'L', 3, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (81, 6, 'aud_component_hijack', 'Intent组件隐式调用风险', 'L', 2, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (82, 6, 'aud_reflect', '反射调用风险', 'L', 3, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (83, 6, 'aud_user_privacy_info', '用户隐私信息检测', 'M', 4, 0, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (84, 6, 'aud_override_perm', '覆盖权限验证', 'L', 0, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (85, 6, 'aud_only_write_perm', '仅定义Provider writePermission', 'L', 0, 0, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (86, 6, 'aud_perm_check', '权限检查', 'L', 0, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (87, 6, 'aud_network_conf', '缺少网络安全配置', 'M', 0, 0, 11, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (88, 6, 'aud_miss_receiver_perm', '缺少Receiver权限', 'M', 0, 0, 12, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (89, 6, 'aud_miss_broadcaster_perm', '缺少Broadcaster权限', 'L', 0, 0, 13, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (90, 6, 'aud_miss_component_perm', '缺少导出的标志或组件权限', 'M', 0, 0, 14, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (91, 6, 'aud_provider_config', 'ContentProvider访问路径配置安全', 'M', 0, 0, 15, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (92, 6, 'aud_sticky_broadcast', '粘性广播使用风险', 'L', 0, 0, 16, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (93, 6, 'aud_coms_activity', 'Activity组件导出风险', 'M', 4, 0, 17, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (94, 6, 'aud_coms_service', 'Service组件导出风险', 'M', 4, 0, 18, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (95, 6, 'aud_coms_receiver', 'Broadcast Receiver组件导出风险', 'M', 4, 0, 19, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (96, 6, 'aud_coms_provider', 'Content Provider组件导出风险', 'M', 4, 0, 20, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (97, 6, 'cve_gif_drawable', 'Android-gif-Drawable远程代码执行漏洞', 'H', 6, 0, 21, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (98, 6, 'cvs_fragment_risk', 'Fragment注入攻击漏洞', 'M', 4, 0, 22, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (99, 6, 'cvs_intent_risk', 'Intent Scheme URL攻击漏洞', 'M', 4, 0, 23, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (100, 8, 'aud_webview_fileurl', '“应用克隆”漏洞攻击风险', 'H', 6, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (101, 8, 'aud_inject', '动态注入攻击风险', 'H', 6, 1, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (102, 8, 'cvs_wv_inject', 'Webview远程代码执行漏洞', 'H', 5, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (103, 8, 'aud_webview_hide_interface', '未移除有风险的Webview系统隐藏接口漏洞', 'H', 5, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (104, 8, 'aud_unzip', 'zip文件解压目录遍历漏洞', 'M', 4, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (105, 8, 'cvs_anydown', '下载任意apk漏洞', 'H', 5, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (106, 8, 'cvs_refuse_service', '拒绝服务攻击漏洞', 'M', 4, 1, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (107, 8, 'aud_sdcard_loaddex', '从sdcard加载dex风险', 'M', 4, 0, 8, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (108, 8, 'aud_sdcard_loadso', '从sdcard加载so风险', 'M', 4, 0, 9, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (109, 8, 'aud_stack_protect', '未使用编译器堆栈保护技术风险', 'L', 3, 0, 10, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (110, 8, 'aud_random_space', '未使用地址空间随机化技术风险', 'L', 3, 0, 11, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (111, 8, 'cvs_emulator_run', '模拟器运行风险', 'L', 3, 1, 12, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (112, 8, 'aud_root_device', 'Root设备运行风险', 'L', 3, 1, 13, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (113, 8, 'aud_risk_webBrowser', '不安全的浏览器调用漏洞', 'M', 4, 0, 14, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (114, 8, 'aud_parasitic_push', '“寄生推”云控风险检测', 'L', 2, 0, 15, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (115, 8, 'cvs_run_other_program', '运行其他可执行程序漏洞', 'M', 4, 0, 16, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (116, 8, 'cvs_upnp', 'libupnp缓冲区溢出漏洞', 'M', 0, 0, 17, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (117, 8, 'aud_intent_redirect', 'Intent重定向', 'M', 0, 0, 18, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (118, 9, 'h5_storage', 'Web Storage数据泄露风险', 'L', 2, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (119, 9, 'h5_websql', 'WebSQL注入漏洞', 'H', 5, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (120, 9, 'h5_innerhtml', 'InnerHTML的XSS攻击漏洞', 'H', 5, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (121, 10, 'sec_other_sdk', '第三方SDK检测', 'I', 0, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (122, 11, 'sec_words', '敏感词信息', 'I', 0, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (123, 12, 'cvs_global_exception', '全局异常检测', 'L', 4, 0, 1, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (124, 12, 'aud_ip_address', 'IP地址检测', 'I', 0, 0, 2, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (125, 12, 'aud_state_crypt_algorithm', '国密算法检测', 'I', 0, 0, 3, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (126, 12, 'cvs_sharedprefs_commit', 'SharedPreferences使用commit提交数据检测', 'I', 0, 0, 4, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (127, 12, 'aud_start_behavior', '自启行为检测', 'L', 4, 0, 5, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (128, 12, 'aud_state_crypt_algorithm_sm2', 'SM2国密算法检测', 'I', 0, 0, 6, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (129, 12, 'aud_state_crypt_algorithm_sm3', 'SM3国密算法检测', 'I', 0, 0, 7, 1, '');
INSERT INTO `ceping_ad_audit_item` VALUES (130, 12, 'aud_state_crypt_algorithm_sm4', 'SM4国密算法检测', 'I', 0, 0, 8, 1, '');

SET FOREIGN_KEY_CHECKS = 1;
