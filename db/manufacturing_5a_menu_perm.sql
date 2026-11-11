-- ============================================================
-- VentureMES 阶段5a：MES 工艺主数据 菜单 + API + Casbin
-- 角色：admin (authority_id=888)
-- ============================================================

USE `gva`;

-- ============================================================
-- 一、清理旧数据
-- ============================================================
DELETE FROM `sys_authority_menus`
  WHERE `sys_base_menu_id` IN (SELECT id FROM (
    SELECT id FROM `sys_base_menus` WHERE `name` IN ('mfgWorkstation','mfgTeam','mfgRouting')
  ) t);
DELETE FROM `sys_base_menus` WHERE `name` IN ('mfgWorkstation','mfgTeam','mfgRouting');

DELETE FROM `casbin_rule` WHERE `v1` LIKE '/mfg/workstation%' OR `v1` LIKE '/mfg/team%' OR `v1` LIKE '/mfg/routing%';
DELETE FROM `sys_apis`    WHERE `path` LIKE '/mfg/workstation%' OR `path` LIKE '/mfg/team%' OR `path` LIKE '/mfg/routing%';

-- ============================================================
-- 二、子菜单（挂在已有的 manufacturing 父菜单下）
-- ============================================================
SET @pid = (SELECT id FROM `sys_base_menus` WHERE `name`='manufacturing' LIMIT 1);

INSERT INTO `sys_base_menus`
  (`created_at`,`updated_at`,`menu_level`,`parent_id`,`path`,`name`,`hidden`,`component`,`sort`,`keep_alive`,`default_menu`,`title`,`icon`)
VALUES
  (NOW(),NOW(),0,@pid,'mfgWorkstation','mfgWorkstation',0,'view/manufacturing/workstation/index.vue',5,0,0,'工位管理','monitor'),
  (NOW(),NOW(),0,@pid,'mfgTeam',       'mfgTeam',       0,'view/manufacturing/team/index.vue',       6,0,0,'班组管理','user-filled'),
  (NOW(),NOW(),0,@pid,'mfgRouting',    'mfgRouting',    0,'view/manufacturing/routing/index.vue',    7,0,0,'工艺路线','set-up');

-- ============================================================
-- 三、admin 授权
-- ============================================================
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`,`sys_base_menu_id`)
  SELECT 888, id FROM `sys_base_menus` WHERE `name` IN ('mfgWorkstation','mfgTeam','mfgRouting');

-- ============================================================
-- 四、注册 sys_apis
-- ============================================================
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  -- 工位
  (NOW(),NOW(),'/mfg/workstation',      '创建工位',  'MES-工艺主数据','POST'),
  (NOW(),NOW(),'/mfg/workstation',      '更新工位',  'MES-工艺主数据','PUT'),
  (NOW(),NOW(),'/mfg/workstation',      '删除工位',  'MES-工艺主数据','DELETE'),
  (NOW(),NOW(),'/mfg/workstation',      '查询单个工位','MES-工艺主数据','GET'),
  (NOW(),NOW(),'/mfg/workstation/list', '分页查询工位','MES-工艺主数据','GET'),
  (NOW(),NOW(),'/mfg/workstation/all',  '获取全部启用工位','MES-工艺主数据','GET'),
  -- 班组
  (NOW(),NOW(),'/mfg/team',             '创建班组',  'MES-工艺主数据','POST'),
  (NOW(),NOW(),'/mfg/team',             '更新班组',  'MES-工艺主数据','PUT'),
  (NOW(),NOW(),'/mfg/team',             '删除班组',  'MES-工艺主数据','DELETE'),
  (NOW(),NOW(),'/mfg/team',             '查询单个班组','MES-工艺主数据','GET'),
  (NOW(),NOW(),'/mfg/team/list',        '分页查询班组','MES-工艺主数据','GET'),
  -- 工艺路线
  (NOW(),NOW(),'/mfg/routing',          '创建工艺路线','MES-工艺主数据','POST'),
  (NOW(),NOW(),'/mfg/routing',          '更新工艺路线','MES-工艺主数据','PUT'),
  (NOW(),NOW(),'/mfg/routing',          '删除工艺路线','MES-工艺主数据','DELETE'),
  (NOW(),NOW(),'/mfg/routing/status',   '启用/停用工艺路线','MES-工艺主数据','POST'),
  (NOW(),NOW(),'/mfg/routing',          '查询单个工艺路线','MES-工艺主数据','GET'),
  (NOW(),NOW(),'/mfg/routing/list',     '分页查询工艺路线','MES-工艺主数据','GET'),
  (NOW(),NOW(),'/mfg/routing/activeByProduct','按成品取启用工艺路线','MES-工艺主数据','GET');

-- ============================================================
-- 五、Casbin 授权
-- ============================================================
INSERT INTO `casbin_rule` (`ptype`,`v0`,`v1`,`v2`)
  SELECT 'p', '888', `path`, `method` FROM `sys_apis`
  WHERE `path` LIKE '/mfg/workstation%' OR `path` LIKE '/mfg/team%' OR `path` LIKE '/mfg/routing%';

-- ============================================================
-- 完成。重启后端 + 重新登录前端。
-- ============================================================
