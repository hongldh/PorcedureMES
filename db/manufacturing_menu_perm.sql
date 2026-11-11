-- ============================================================
-- VentureMES 阶段4：MES 模块 菜单 + API + Casbin 权限注入
-- 适用：mfg_bom + mfg_workOrder + mfg_issue + mfg_receipt
-- 角色：admin (authority_id=888)
-- ============================================================

USE `gva`;

-- ============================================================
-- 一、清理旧数据（幂等）
-- ============================================================
DELETE FROM `sys_authority_menus`
  WHERE `sys_base_menu_id` IN (SELECT id FROM (
    SELECT id FROM `sys_base_menus` WHERE `name` LIKE 'mfg%' OR `name`='manufacturing'
  ) t);
DELETE FROM `sys_base_menus` WHERE `name` LIKE 'mfg%' OR `name`='manufacturing';

DELETE FROM `casbin_rule` WHERE `v1` LIKE '/mfg/%';
DELETE FROM `sys_apis` WHERE `path` LIKE '/mfg/%';

-- ============================================================
-- 二、一级菜单
-- ============================================================
INSERT INTO `sys_base_menus`
  (`created_at`,`updated_at`,`menu_level`,`parent_id`,`path`,`name`,`hidden`,`component`,`sort`,`keep_alive`,`default_menu`,`title`,`icon`)
VALUES
  (NOW(),NOW(),0,0,'manufacturing','manufacturing',0,'view/routerHolder.vue',200,0,0,'制造执行','setting');

SET @pid = (SELECT id FROM `sys_base_menus` WHERE `name`='manufacturing' LIMIT 1);

-- ============================================================
-- 三、子菜单 4 项
-- ============================================================
INSERT INTO `sys_base_menus`
  (`created_at`,`updated_at`,`menu_level`,`parent_id`,`path`,`name`,`hidden`,`component`,`sort`,`keep_alive`,`default_menu`,`title`,`icon`)
VALUES
  (NOW(),NOW(),0,@pid,'mfgBom',       'mfgBom',       0,'view/manufacturing/bom/index.vue',        1,0,0,'BOM 物料清单','document'),
  (NOW(),NOW(),0,@pid,'mfgWorkOrder', 'mfgWorkOrder', 0,'view/manufacturing/workOrder/index.vue',  2,0,0,'生产工单',   'tickets'),
  (NOW(),NOW(),0,@pid,'mfgIssue',     'mfgIssue',     0,'view/manufacturing/mfgIssue/index.vue',   3,0,0,'领料单',     'shopping-cart'),
  (NOW(),NOW(),0,@pid,'mfgReceipt',   'mfgReceipt',   0,'view/manufacturing/mfgReceipt/index.vue', 4,0,0,'完工入库',   'box');

-- ============================================================
-- 四、给 admin 授权
-- ============================================================
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`,`sys_base_menu_id`)
  SELECT 888, id FROM `sys_base_menus`
  WHERE `name` IN ('manufacturing') OR `name` LIKE 'mfg%';

-- ============================================================
-- 五、注册 sys_apis
-- ============================================================
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  -- BOM
  (NOW(),NOW(),'/mfg/bom',                  '创建BOM',         'MES-BOM','POST'),
  (NOW(),NOW(),'/mfg/bom',                  '更新BOM',         'MES-BOM','PUT'),
  (NOW(),NOW(),'/mfg/bom',                  '删除BOM',         'MES-BOM','DELETE'),
  (NOW(),NOW(),'/mfg/bom/status',           '启用/停用BOM',    'MES-BOM','POST'),
  (NOW(),NOW(),'/mfg/bom',                  '查询单个BOM',     'MES-BOM','GET'),
  (NOW(),NOW(),'/mfg/bom/list',             '分页查询BOM',     'MES-BOM','GET'),
  (NOW(),NOW(),'/mfg/bom/activeByProduct',  '按成品取启用BOM', 'MES-BOM','GET'),
  -- 工单
  (NOW(),NOW(),'/mfg/workOrder',            '创建工单',  'MES-工单','POST'),
  (NOW(),NOW(),'/mfg/workOrder',            '更新工单',  'MES-工单','PUT'),
  (NOW(),NOW(),'/mfg/workOrder',            '删除工单',  'MES-工单','DELETE'),
  (NOW(),NOW(),'/mfg/workOrder/release',    '下达工单',  'MES-工单','POST'),
  (NOW(),NOW(),'/mfg/workOrder/close',      '关闭工单',  'MES-工单','POST'),
  (NOW(),NOW(),'/mfg/workOrder/void',       '作废工单',  'MES-工单','POST'),
  (NOW(),NOW(),'/mfg/workOrder',            '查询单个工单','MES-工单','GET'),
  (NOW(),NOW(),'/mfg/workOrder/list',       '分页查询工单','MES-工单','GET'),
  -- 领料
  (NOW(),NOW(),'/mfg/issue',                '创建领料单',  'MES-领料','POST'),
  (NOW(),NOW(),'/mfg/issue',                '更新领料单',  'MES-领料','PUT'),
  (NOW(),NOW(),'/mfg/issue',                '删除领料单',  'MES-领料','DELETE'),
  (NOW(),NOW(),'/mfg/issue/confirm',        '确认领料',    'MES-领料','POST'),
  (NOW(),NOW(),'/mfg/issue/void',           '作废领料单',  'MES-领料','POST'),
  (NOW(),NOW(),'/mfg/issue',                '查询单个领料单','MES-领料','GET'),
  (NOW(),NOW(),'/mfg/issue/list',           '分页查询领料单','MES-领料','GET'),
  -- 完工入库
  (NOW(),NOW(),'/mfg/receipt',              '创建完工单',  'MES-完工','POST'),
  (NOW(),NOW(),'/mfg/receipt',              '更新完工单',  'MES-完工','PUT'),
  (NOW(),NOW(),'/mfg/receipt',              '删除完工单',  'MES-完工','DELETE'),
  (NOW(),NOW(),'/mfg/receipt/confirm',      '确认完工入库','MES-完工','POST'),
  (NOW(),NOW(),'/mfg/receipt/void',         '作废完工单',  'MES-完工','POST'),
  (NOW(),NOW(),'/mfg/receipt',              '查询单个完工单','MES-完工','GET'),
  (NOW(),NOW(),'/mfg/receipt/list',         '分页查询完工单','MES-完工','GET');

-- ============================================================
-- 六、Casbin 授权 admin 全部 /mfg/*
-- ============================================================
INSERT INTO `casbin_rule` (`ptype`,`v0`,`v1`,`v2`)
  SELECT 'p', '888', `path`, `method` FROM `sys_apis` WHERE `path` LIKE '/mfg/%';

-- ============================================================
-- 完成。请重新登录前端以加载新菜单与权限缓存。
-- ============================================================
