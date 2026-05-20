-- ============================================================
-- VentureMES 进销存模块 菜单 + API + Casbin 权限注入脚本
-- 数据库：gva
-- 适用范围：阶段 1 主数据(7) + 阶段 2 单据(6) + 库存查询(2) = 15 个页面
-- 角色：admin (authority_id=888)
--
-- 使用方法：
--   mysql -uroot -p gva < db/inventory_menu_perm.sql
-- 或在 Navicat 中直接执行本脚本。
--
-- 幂等：脚本可重复执行，先 DELETE 旧数据再 INSERT。
-- ============================================================

USE `gva`;

-- ============================================================
-- 一、清理旧数据（幂等）
-- ============================================================

-- 1) 清理菜单（按 path 前缀识别本模块）
DELETE FROM `sys_authority_menus`
  WHERE `sys_base_menu_id` IN (SELECT id FROM (
    SELECT id FROM `sys_base_menus` WHERE `name` LIKE 'inv%'
  ) t);
DELETE FROM `sys_base_menus` WHERE `name` LIKE 'inv%';

-- 2) 清理 API
DELETE FROM `casbin_rule` WHERE `v1` LIKE '/inv/%';
DELETE FROM `sys_apis` WHERE `path` LIKE '/inv/%';

-- ============================================================
-- 二、插入一级菜单（父级"进销存管理"）
-- ============================================================
INSERT INTO `sys_base_menus`
  (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`)
VALUES
  (NOW(), NOW(), 0, 0, 'inventory', 'inventory', 0, 'view/routerHolder.vue', 100, 0, 0, '进销存管理', 'goods');

-- 取出父菜单 ID
SET @pid = (SELECT id FROM `sys_base_menus` WHERE `name`='inventory' LIMIT 1);

-- ============================================================
-- 三、插入子菜单（15 个页面）
--    component 路径必须与 web/src/view/inventory/<dir>/index.vue 对齐
-- ============================================================
INSERT INTO `sys_base_menus`
  (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`)
VALUES
  -- 主数据（阶段 1）
  (NOW(), NOW(), 0, @pid, 'invCategory',  'invCategory',   0, 'view/inventory/category/index.vue',   1, 0, 0, '商品分类', 'menu'),
  (NOW(), NOW(), 0, @pid, 'invUnit',      'invUnit',       0, 'view/inventory/unit/index.vue',       2, 0, 0, '计量单位', 'scale-to-original'),
  (NOW(), NOW(), 0, @pid, 'invProduct',   'invProduct',    0, 'view/inventory/product/index.vue',    3, 0, 0, '商品资料', 'goods-filled'),
  (NOW(), NOW(), 0, @pid, 'invWarehouse', 'invWarehouse',  0, 'view/inventory/warehouse/index.vue',  4, 0, 0, '仓库管理', 'house'),
  (NOW(), NOW(), 0, @pid, 'invLocation',  'invLocation',   0, 'view/inventory/location/index.vue',   5, 0, 0, '货位管理', 'location'),
  (NOW(), NOW(), 0, @pid, 'invSupplier',  'invSupplier',   0, 'view/inventory/supplier/index.vue',   6, 0, 0, '供应商',   'avatar'),
  (NOW(), NOW(), 0, @pid, 'invCustomer',  'invInvCustomer',0, 'view/inventory/customer/index.vue',   7, 0, 0, '客户管理', 'user'),
  -- 业务单据（阶段 2）
  (NOW(), NOW(), 0, @pid, 'invPurchaseOrder',  'invPurchaseOrder',  0, 'view/inventory/purchaseOrder/index.vue',  10, 0, 0, '采购入库', 'shopping-cart'),
  (NOW(), NOW(), 0, @pid, 'invPurchaseReturn', 'invPurchaseReturn', 0, 'view/inventory/purchaseReturn/index.vue', 11, 0, 0, '采购退货', 'refresh-left'),
  (NOW(), NOW(), 0, @pid, 'invSalesOrder',     'invSalesOrder',     0, 'view/inventory/salesOrder/index.vue',     12, 0, 0, '销售出库', 'sell'),
  (NOW(), NOW(), 0, @pid, 'invSalesReturn',    'invSalesReturn',    0, 'view/inventory/salesReturn/index.vue',    13, 0, 0, '销售退货', 'refresh-right'),
  (NOW(), NOW(), 0, @pid, 'invTransfer',       'invTransfer',       0, 'view/inventory/transfer/index.vue',       14, 0, 0, '仓间调拨', 'right'),
  (NOW(), NOW(), 0, @pid, 'invStocktake',      'invStocktake',      0, 'view/inventory/stocktake/index.vue',      15, 0, 0, '库存盘点', 'check'),
  -- 库存查询（阶段 3）
  (NOW(), NOW(), 0, @pid, 'invStockBalance',   'invStockBalance',   0, 'view/inventory/stockBalance/index.vue',   20, 0, 0, '库存余额', 'data-analysis'),
  (NOW(), NOW(), 0, @pid, 'invStockFlow',      'invStockFlow',      0, 'view/inventory/stockFlow/index.vue',      21, 0, 0, '库存流水', 'tickets');

-- ============================================================
-- 四、给 admin (authority_id=888) 授权所有上述菜单
-- ============================================================
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`)
  SELECT 888, id FROM `sys_base_menus`
  WHERE `name` IN ('inventory') OR `name` LIKE 'inv%';

-- ============================================================
-- 五、注册 sys_apis（66 条：主数据 + 单据 + 库存查询）
-- ============================================================

-- 主数据：每张表 7 端点 (POST/PUT/DELETE/DELETE batch/GET/GET list/GET all)
-- 共 7 张表 × 7 = 49 条
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  (NOW(),NOW(),'/inv/category','创建商品分类','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/category','更新商品分类','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/category','删除商品分类','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/category/batch','批量删除商品分类','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/category','查询单个商品分类','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/category/list','分页查询商品分类','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/category/all','获取全部商品分类','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/unit','创建单位','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/unit','更新单位','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/unit','删除单位','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/unit/batch','批量删除单位','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/unit','查询单个单位','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/unit/list','分页查询单位','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/unit/all','获取全部单位','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/product','创建商品','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/product','更新商品','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/product','删除商品','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/product/batch','批量删除商品','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/product','查询单个商品','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/product/list','分页查询商品','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/product/all','获取全部商品','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/warehouse','创建仓库','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/warehouse','更新仓库','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/warehouse','删除仓库','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/warehouse/batch','批量删除仓库','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/warehouse','查询单个仓库','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/warehouse/list','分页查询仓库','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/warehouse/all','获取全部仓库','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/location','创建货位','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/location','更新货位','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/location','删除货位','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/location/batch','批量删除货位','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/location','查询单个货位','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/location/list','分页查询货位','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/location/byWarehouse','按仓库取货位','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/supplier','创建供应商','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/supplier','更新供应商','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/supplier','删除供应商','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/supplier/batch','批量删除供应商','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/supplier','查询单个供应商','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/supplier/list','分页查询供应商','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/supplier/all','获取全部供应商','进销存-主数据','GET'),

  (NOW(),NOW(),'/inv/customer','创建客户','进销存-主数据','POST'),
  (NOW(),NOW(),'/inv/customer','更新客户','进销存-主数据','PUT'),
  (NOW(),NOW(),'/inv/customer','删除客户','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/customer/batch','批量删除客户','进销存-主数据','DELETE'),
  (NOW(),NOW(),'/inv/customer','查询单个客户','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/customer/list','分页查询客户','进销存-主数据','GET'),
  (NOW(),NOW(),'/inv/customer/all','获取全部客户','进销存-主数据','GET');

-- 业务单据：每个 7 端点 (POST/PUT/DELETE/POST confirm/POST void/GET/GET list)
-- 共 6 单据 × 7 = 42 条
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  (NOW(),NOW(),'/inv/purchaseOrder',         '创建采购入库单',  '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseOrder',         '更新采购入库单',  '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/purchaseOrder',         '删除采购入库单草稿','进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/purchaseOrder/confirm', '确认采购入库',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseOrder/void',    '作废采购入库单',   '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseOrder',         '查询单个采购入库单','进销存-单据','GET'),
  (NOW(),NOW(),'/inv/purchaseOrder/list',    '分页采购入库单',   '进销存-单据','GET'),

  (NOW(),NOW(),'/inv/purchaseReturn',         '创建采购退货单',  '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseReturn',         '更新采购退货单',  '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/purchaseReturn',         '删除采购退货单草稿','进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/purchaseReturn/confirm', '确认采购退货',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseReturn/void',    '作废采购退货单',   '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/purchaseReturn',         '查询单个采购退货单','进销存-单据','GET'),
  (NOW(),NOW(),'/inv/purchaseReturn/list',    '分页采购退货单',   '进销存-单据','GET'),

  (NOW(),NOW(),'/inv/salesOrder',         '创建销售出库单',  '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesOrder',         '更新销售出库单',  '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/salesOrder',         '删除销售出库单草稿','进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/salesOrder/confirm', '确认销售出库',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesOrder/void',    '作废销售出库单',   '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesOrder',         '查询单个销售出库单','进销存-单据','GET'),
  (NOW(),NOW(),'/inv/salesOrder/list',    '分页销售出库单',   '进销存-单据','GET'),

  (NOW(),NOW(),'/inv/salesReturn',         '创建销售退货单',  '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesReturn',         '更新销售退货单',  '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/salesReturn',         '删除销售退货单草稿','进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/salesReturn/confirm', '确认销售退货',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesReturn/void',    '作废销售退货单',   '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/salesReturn',         '查询单个销售退货单','进销存-单据','GET'),
  (NOW(),NOW(),'/inv/salesReturn/list',    '分页销售退货单',   '进销存-单据','GET'),

  (NOW(),NOW(),'/inv/transfer',         '创建调拨单',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/transfer',         '更新调拨单',     '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/transfer',         '删除调拨单草稿', '进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/transfer/confirm', '确认调拨',       '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/transfer/void',    '作废调拨单',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/transfer',         '查询单个调拨单', '进销存-单据','GET'),
  (NOW(),NOW(),'/inv/transfer/list',    '分页调拨单',     '进销存-单据','GET'),

  (NOW(),NOW(),'/inv/stocktake',         '创建盘点单',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/stocktake',         '更新盘点单',     '进销存-单据','PUT'),
  (NOW(),NOW(),'/inv/stocktake',         '删除盘点单草稿', '进销存-单据','DELETE'),
  (NOW(),NOW(),'/inv/stocktake/confirm', '确认盘点',       '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/stocktake/void',    '作废盘点单',     '进销存-单据','POST'),
  (NOW(),NOW(),'/inv/stocktake',         '查询单个盘点单', '进销存-单据','GET'),
  (NOW(),NOW(),'/inv/stocktake/list',    '分页盘点单',     '进销存-单据','GET');

-- 库存查询：2 条
INSERT INTO `sys_apis` (`created_at`,`updated_at`,`path`,`description`,`api_group`,`method`) VALUES
  (NOW(),NOW(),'/inv/stock/balance', '库存余额查询', '进销存-库存','GET'),
  (NOW(),NOW(),'/inv/stock/flow',    '库存流水查询', '进销存-库存','GET');

-- ============================================================
-- 六、Casbin 授权（admin = authority_id 888）
--   表结构：(ptype, v0=role, v1=path, v2=method)
-- ============================================================
INSERT INTO `casbin_rule` (`ptype`,`v0`,`v1`,`v2`)
  SELECT 'p', '888', `path`, `method` FROM `sys_apis` WHERE `path` LIKE '/inv/%';

-- ============================================================
-- 完成。请重新登录前端以加载新菜单与权限缓存。
-- ============================================================
