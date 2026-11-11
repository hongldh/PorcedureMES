# VentureMES 交接文档

> 基于 gin-vue-admin（Go + Vue3）改造的中小制造业「进销存 + MES」一体化系统。
> 截止阶段 5b。不含财务 / WMS 货位作业 / 多渠道 OMS / 质检 / SPC。

## 1. 项目功能

| 模块 | 功能 | 状态 |
|------|------|------|
| 主数据 | 商品 / 分类 / 单位 / 仓库 / 货位 / 供应商 / 客户 | ✅ |
| 进销存 | 采购入库、采购退货、销售出库、销售退货、调拨、盘点；移动加权平均成本；库存余额 / 流水查询 | ✅ |
| MES 基础 | BOM（多版本）、生产工单、领料出库、完工入库（物料齐套校验，成本 = 领料金额 / 完工数） | ✅ |
| MES 工艺 | 工位、班组（成员）、工艺路线（多版本、工序） | ✅ |
| MES 执行 | 工单下达展开工序快照、工序派工（员工/班组/工位）、员工报工、作废依赖链校验、完工入库受末道工序合格数约束 | ✅ |
| 批次追溯 | — | 🔜 5d |
| 报表看板 | — | 🔜 阶段 6 |

核心规则：所有库存变动经 `StockEngine.Apply` 写 `inv_stock_flow` 后更新 `inv_stock`；单据 3 态 草稿(1)/已执行(3)/作废(4)，作废 = 反向冲销；单号 `{前缀}{YYYYMMDD}{3位序}`。

## 2. 组成模块

```
server/                              Go 后端（GVA 五层：model → service → api → router → initialize 注册）
├── model|service|api/v1|router/
│   ├── inventory/                   进销存（inv_* 表，93 端点）
│   └── manufacturing/               MES（mfg_* 表，53 端点）
├── initialize/gorm_biz.go           业务表 AutoMigrate
└── initialize/router_biz.go         业务路由挂载
web/src/
├── api/{inventory,manufacturing}/   接口封装
└── view/{inventory,manufacturing}/  页面（进销存 15 页；MES 8 页）
db/*.sql                             菜单 + sys_apis + Casbin 授权脚本（幂等，admin=888）
docs/learning/*.md                   说明文档
```

GVA 原有代码改 5 个注册点：`api/v1/enter.go`、`service/enter.go`、`router/enter.go`、`gorm_biz.go`、`router_biz.go`。

## 3. 技术特点

1. **分层架构清晰**：Model → Service → API → Router 依赖单向向下，Service 不依赖 Gin，便于测试与替换传输层。
2. **事务边界内聚在 Service**：Confirm/Void 以 `db.Transaction` 包裹跨表写入，API 层无事务逻辑。
3. **悲观锁 + 幂等**：库存行、工序行、单号序列 `FOR UPDATE`；状态机校验保证重复请求无副作用。
4. **Append-only 事件溯源**：流水、报工只追加不修改，余额为派生状态，可审计、可重建。
5. **统一领域引擎**：`StockEngine.Apply` 抽象"方向+数量+成本"，新增单据类型只调用不重写，符合开闭原则。
6. **快照隔离主数据**：工单持有 BOM/工艺副本，后续修改不影响在制工单。
7. **约束前移至领域层**：库存≥0、工序顺序、齐套、权属等不变量在 Service 强制，前端仅做 UX。
8. **组合优于继承**：`GVA_MODEL` 嵌入复用字段，`ServiceGroup` 嵌入聚合单例，路由接口断言 + `bind()` 泛化注册。
9. **前端 Composition API 复用**：逻辑抽 `useDocPage`，UI 抽通用组件，单据页仅配置注入。
10. **对框架零侵入**：业务独立成包，仅 5 个注册点接入，框架升级冲突面极小。
11. **精度控制**：数量/金额 `decimal(18,4)`，比较统一用 `qtyEps` 容差。
12. **幂等部署脚本**：权限 SQL 先删后插，Casbin 由 `sys_apis` 派生，可重放。

## 4. 部署步骤

环境：Go ≥ 1.22、Node ≥ 18、MySQL ≥ 5.7（需支持严格模式）、Redis 可选。

```powershell
# ① 后端（首次启动走 GVA /init 向导建库；之后 AutoMigrate 自动建业务表）
cd server
go run .                             # 默认 :8888

# ② 前端
cd web
npm install
npm run dev                          # 默认 :8080，admin / 123456

# ③ 业务菜单与权限（按顺序执行一次，Navicat 或 cmd）
cmd /c "mysql -uroot -p<密码> gva < db\inventory_menu_perm.sql"
cmd /c "mysql -uroot -p<密码> gva < db\manufacturing_menu_perm.sql"
cmd /c "mysql -uroot -p<密码> gva < db\manufacturing_5a_menu_perm.sql"
cmd /c "mysql -uroot -p<密码> gva < db\manufacturing_5b_menu_perm.sql"

# ④ 重启后端（重载 Casbin）→ 重新登录前端
```

注意：PowerShell 不能直接用 `<` 重定向，需走 `cmd /c`；每次新增 API 必须补 `sys_apis` + `casbin_rule` 并重启后端。
