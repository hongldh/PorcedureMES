-- ============================================================
-- VentureMES 阶段5b：MES 派工 + 报工 菜单 + API + Casbin
-- 角色：admin (authority_id=888)
-- 依赖：已执行 manufacturing_menu_perm.sql（存在父菜单 manufacturing）
-- ============================================================

USE `gva`;

-- ============================================================
-- 一、清理旧数据（幂等）
-- ============================================================
DELETE FROM `sys_authority_menus`
  WHERE `sys_base_menu_id` IN (SELECT id FROM (
    SELECT id FROM `sys_base_menus` WHERE `name` = 'mfgWorkReport'
  ) t);
DELETE FROM `sys_base_menus` WHERE `name` = 'mfgWorkReport';

DELETE FROM `casbin_rule`
  WHERE `v1` IN ('/mfg/workOrder/ops', '/mfg/op/dispatch')
     OR `v1` LIKE '/mfg/report%';
DELETE FROM `sys_apis`
  WHERE `path` IN ('/mfg/workOrder/ops', '/mfg/op/dispatch')
     OR `path` LIKE '/mfg/report%';

-- ============================================================
-- 二、子菜单：报工（挂在 manufacturing 父菜单下）
--    派工入口集成在"生产工单 → 查看"对话框内，不单独建菜单
-- ============================================================
SET @pid = (SELECT id FROM `sys_base_menus` WHERE `name`='manufacturing' LIMIT 1);

INSERT INTO `sys_base_menus`
  (`created_at`,`updated_at`,`menu_level`,`parent_id`,`path`,`name`,`hidden`,`component`,`sort`,`keep_alive`,`default_menu`,`title`,`icon`)
VALUES
  (NOW(),NOW(),0,@pid,'mfgWorkReport','mfgWorkReport',0,'view/manufacturing/workReport/index.vue',8,0,0,'工序报工','edit-pen');

-- ============================================================
-- 三、admin 授权菜单
-- ============================================================
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`,`sys_base_menu_id`)
  SELECT 888, id FROM `sys_base_menus` WHERE `name` = 'mfgWorkReport';

-- ============================================================
-- 四、注册 sys_apis（6 个端点）
-- ============================================================
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  (NOW(),NOW(),'/mfg/workOrder/ops', '查询工单工序列表', 'MES-派工报工','GET'),
  (NOW(),NOW(),'/mfg/op/dispatch',   '工序派工/改派',     'MES-派工报工','POST'),
  (NOW(),NOW(),'/mfg/report',        '提交报工',          'MES-派工报工','POST'),
  (NOW(),NOW(),'/mfg/report/void',   '作废报工',          'MES-派工报工','POST'),
  (NOW(),NOW(),'/mfg/report/list',   '报工记录分页',      'MES-派工报工','GET'),
  (NOW(),NOW(),'/mfg/report/myOps',  '我的待报工工序',    'MES-派工报工','GET');

-- ============================================================
-- 五、Casbin 授权 admin
-- ============================================================
INSERT INTO `casbin_rule` (`ptype`,`v0`,`v1`,`v2`)
  SELECT 'p', '888', `path`, `method` FROM `sys_apis`
  WHERE `path` IN ('/mfg/workOrder/ops', '/mfg/op/dispatch')
     OR `path` LIKE '/mfg/report%';

-- ============================================================
-- 完成。重启后端（建 mfg_work_order_op / mfg_work_report 两表）+ 重新登录前端。
-- ============================================================
