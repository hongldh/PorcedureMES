# 08 - 阶段 5b：MES 派工 + 报工（工序级车间执行）

## 一、本子阶段定位

阶段 4 的工单是"黑盒"：下达 → 领料 → 完工入库，中间发生了什么系统不知道。
5a 补了"工序怎么定义"（工艺路线），5b 把工序**落到工单上并跑起来**：

- 工单下达时按工艺路线**展开工序快照**
- 把每道工序**派给员工 / 班组 / 工位**
- 员工按工序**报工**（合格数 / 不良数 / 工时）
- 完工入库数量**受末道工序合格数约束**

完成 5b 后，系统第一次能回答："这批货现在到哪道工序了？谁在做？做了多少？"——这是 MES 区别于进销存的核心能力。

---

## 二、决策回顾（沿用阶段 5 开工时的拍板）

| 维度 | 选择 | 在 5b 的体现 |
|------|------|-------------|
| 工序组织 | 工艺路线模式 | 下达时从 `mfg_routing_item` 复制到 `mfg_work_order_op` |
| 报工 | 员工独立报工 | 报工记录带 `UserID`；派给个人必须本人报，派给班组必须组员报 |
| 数量校验 | 严格不允许 | 超计划 / 跳工序 / 入库超末道合格 全部拒绝 |
| 质检 / SPC | 暂不做 | `RejectQty` 只记数，不做不良处理流程 |

---

## 三、新增表（2 张）

```
mfg_work_order_op     工单工序 = 工艺路线快照 + 执行态 + 派工信息
mfg_work_report       报工记录（单号前缀 BG，即时生效，无草稿）
```

### 3.1 MfgWorkOrderOp 三段式字段

```
快照段（下达时从 RoutingItem 复制，之后不变）
  Seq / Code / Name / WorkstationID / StdHours / IsKey

执行段（报工累加）
  PlanQty      = 工单计划数
  DoneQty      累计合格
  RejectQty    累计不良
  ActualHours  累计工时

派工段
  AssignedUserID   指派员工（可空）
  TeamID           指派班组（可空）
  Status           1未派工 → 2已派工 → 3进行中 → 4已完工
```

### 3.2 MfgWorkReport 为什么没有草稿态

领料单、完工单都有"草稿 → 已执行"两步，报工却设计成**提交即生效**：

- 车间报工是高频、低价值单笔操作，再加一步"确认"没人会用
- 错了不是"编辑草稿"而是**作废后重报**，和库存流水"只追加不修改"的审计思想一致
- 状态只有 `3 已生效 / 4 作废` 两种

---

## 四、工单状态机扩展

5b 没有新增工单状态，但**状态流转的触发点变了**：

```
草稿(1) ──下达──→ 已下达(2) ──首次报工 或 首次领料──→ 生产中(3)
                                                          │
                                       完工入库累计 ≥ 计划数 ──→ 已完工(4)
```

工序自己也有一个小状态机：

```
未派工(1) ──派工──→ 已派工(2) ──首次报工──→ 进行中(3) ──合格累计 ≥ 计划──→ 已完工(4)
                       ▲                                           │
                       └──────── 作废报工到累计为 0 ────────────────┘（回退到已派工）
```

---

## 五、五条严格校验规则（Service 层核心）

### 5.1 下达必须有工艺路线

5a 把 `RoutingID` 做成可选是为了兼容老工单；5b 起 `Release` 强制：

```go
if d.RoutingID == 0 { return errors.New("下达前必须指定工艺路线") }
```

前端同步把"工艺路线"下拉改为必填。

### 5.2 工序顺序：第 N 道累计 ≤ 第 N-1 道累计合格

```go
var prev MfgWorkOrderOp
tx.Where("work_order_id=? AND seq<?", op.WorkOrderID, op.Seq).Order("seq DESC").First(&prev)
if op.DoneQty+op.RejectQty+reportTotal > prev.DoneQty+qtyEps {
    return fmt.Errorf("上道工序[%s]合格数仅 %.4f ...", prev.Code, prev.DoneQty)
}
```

**为什么用上道"合格"而不是"合格+不良"**：上道判为不良的件，不应该流到下道。这是物理约束，不是业务偏好。

### 5.3 本道累计不超计划

`DoneQty + RejectQty + 本次 ≤ PlanQty`。注意这里用的是合格+不良之和——不良件也消耗了投入。

### 5.4 报工人校验

```go
if op.AssignedUserID > 0 && op.AssignedUserID != userID → 拒绝
if op.TeamID > 0 && 该用户不在 mfg_team_member 中 → 拒绝
```

### 5.5 完工入库 ≤ 末道工序累计合格

`checkLastOpDone` 在完工入库的 `Create` 和 `Confirm` 两处都调：

```go
var last MfgWorkOrderOp
tx.Where("work_order_id=?", id).Order("seq DESC").First(&last)
if totalAfter > last.DoneQty+qtyEps → 拒绝
```

**向后兼容**：查不到工序（5b 之前下达的老工单）直接放行。

---

## 六、作废的"依赖链"校验

作废一条报工不是简单减数，要保证**下游没有依赖它**：

```
作废第 N 道报工
  ├─ 第 N+1 道累计(合格+不良) 不能超过 冲回后的第 N 道合格   → 否则先作废 N+1 道
  └─ 若 N 是末道：工单 ActualQty 不能超过 冲回后的末道合格    → 否则先作废完工入库单
```

这和进销存里"作废采购入库前，必须先作废依赖它的销售出库"是同一种思想：**单据之间的数量依赖是有向图，作废只能从叶子往根走**。

---

## 七、新增 / 修改文件清单

### 后端

```
server/
├── model/manufacturing/
│   └── mfg_execution.go              ★ 新增：MfgWorkOrderOp / MfgWorkReport + 状态常量
├── service/manufacturing/
│   ├── execution.go                  ★ 新增：DispatchService / WorkReportService / checkLastOpDone
│   ├── work_order.go                 ◯ 修改：Release 强制工艺路线 + 展开工序快照；Void 校验报工；新增 GetOps
│   ├── receipt.go                    ◯ 修改：Create / Confirm 加 checkLastOpDone
│   └── enter.go                      ◯ 修改：注册两个新 service
├── service/inventory/sequence.go     ◯ 修改：加 MFG_WORK_REPORT → BG 前缀
├── api/v1/manufacturing/
│   ├── execution.go                  ★ 新增：DispatchApi / WorkReportApi（6 端点）
│   └── enter.go                      ◯ 修改
├── router/manufacturing/
│   ├── execution.go                  ★ 新增：MfgExecutionRouter
│   └── enter.go                      ◯ 修改
├── initialize/
│   ├── gorm_biz.go                   ◯ 修改：AutoMigrate 2 表
│   └── router_biz.go                 ◯ 修改：挂 MfgExecutionRouter
```

### 前端

```
web/src/
├── api/manufacturing/index.js        ◯ 修改：加 dispatchApi / workReportApi
└── view/manufacturing/
    ├── workOrder/index.vue           ◯ 修改：查看对话框加"工序与派工"区块 + 派工子对话框；工艺路线改必填
    └── workReport/index.vue          ★ 新增：我的待报工工序 + 报工记录
```

### SQL + 文档

```
db/manufacturing_5b_menu_perm.sql    ★ 新增：1 子菜单（工序报工）+ 6 API + 6 Casbin
docs/learning/08-阶段5b-MES派工报工.md  ★ 本文
```

### 端点清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /mfg/workOrder/ops?id= | 工单工序列表（含派工信息） |
| POST | /mfg/op/dispatch | 派工 / 改派 |
| POST | /mfg/report | 提交报工（即时生效） |
| POST | /mfg/report/void | 作废报工 |
| GET | /mfg/report/list | 报工记录分页 |
| GET | /mfg/report/myOps | 当前登录人可报工的工序 |

---

## 八、Go 学习要点

### 8.1 `SELECT ... FOR UPDATE` 行锁防并发超报

```go
tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&op, e.OpID)
```

两个员工同时对同一道工序报工，没有行锁会出现：各自读到 `DoneQty=90`，各报 10，都通过 `≤ 100` 校验，最终累计 110。
加了行锁，第二个事务会阻塞到第一个提交后再读，读到 100，被拒绝。

这是 `inventory.StockEngine` 里锁库存余额行的同一招，在 5b 里复用到工序行上。

### 8.2 浮点容差常量化

```go
const qtyEps = 0.0001
if a > b+qtyEps { ... }
```

阶段 4 代码里到处写 `+0.0001`，5b 抽成常量。`decimal(18,4)` 从数据库读出来变 float64 时有舍入误差，不加容差会出现 `99.9999 < 100` 导致"明明做完了却不算完工"。

### 8.3 `gorm.ErrRecordNotFound` 当作正常分支

```go
err := tx.Where(...).Order("seq DESC").First(&prev).Error
if err == nil {
    // 有上道工序，做顺序校验
} else if err != gorm.ErrRecordNotFound {
    return err   // 真正的 DB 错误
}
// ErrRecordNotFound：这是第一道工序，没有上道，跳过
```

"查不到"不是错误而是业务分支的常见模式。注意**不要**用 `if err != nil { return err }` 一把梭。

### 8.4 用 `Pluck` 取单列做 IN 查询

```go
var teamIDs []uint
db.Model(&MfgTeamMember{}).Where("user_id=?", userID).Pluck("team_id", &teamIDs)
db.Where("... team_id IN ?", teamIDs)
```

`Pluck` 比 `Find` 到结构体再遍历取字段省一次内存分配，而且语义清晰。注意 `IN ?` 传空切片会生成 `IN (NULL)`，所以代码里对 `len(teamIDs) > 0` 做了分支。

### 8.5 `Joins` + 表名前缀避免列名歧义

```go
db.Model(&MfgWorkOrderOp{}).
   Joins("JOIN mfg_work_order wo ON wo.id = mfg_work_order_op.work_order_id AND wo.deleted_at IS NULL").
   Where("wo.status IN ?", ...).
   Where("mfg_work_order_op.status IN ?", ...)
```

两张表都有 `status` 列，不带前缀 MySQL 会报 `Column 'status' is ambiguous`。
另外手写 JOIN 时 GORM **不会**自动加软删除条件，`wo.deleted_at IS NULL` 必须自己写。

---

## 九、Vue 学习要点

### 9.1 对话框套对话框：`append-to-body`

派工对话框是从"查看工单"对话框里弹出来的。不加 `append-to-body`，内层 dialog 会被外层的 `overflow: hidden` 裁剪或 z-index 压住。

### 9.2 可选链 `?.` 兜底嵌套 null

后端 Preload 的关联对象在数据缺失时是 `null`（如工序没派工时 `assignedUser: null`）。模板里统一 `row.assignedUser?.nickName || row.team?.name || '-'`，不然渲染直接报错白屏。

### 9.3 前端预校验 + 后端兜底校验

报工对话框里 `合格+不良 ≤ 剩余` 在前端先拦一道，给用户即时反馈；但**后端照样校验**——前端校验是体验，后端校验是安全。

---

## 十、对项目结构的影响

| 维度 | 5a 结束 | 5b 结束 |
|------|---------|---------|
| 数据表 | 34 | **36**（+2） |
| API 端点 | 140 | **146**（+6） |
| 前端页面 | 17 | **18**（+1 报工页；工单页扩展） |
| 单号前缀 | CGRK/CGTH/XSCK/XSTH/DBDD/PDDD/SCGD/LL/RK | **+ BG** |
| 学习文档 | 01–07 | **+ 08** |

---

## 十一、部署步骤

1. **重启后端** `cd server && go run .` — AutoMigrate 建 `mfg_work_order_op`、`mfg_work_report`
2. **执行 SQL** `db/manufacturing_5b_menu_perm.sql`（Navicat 或 `cmd /c "mysql -uroot -p密码 gva < db\manufacturing_5b_menu_perm.sql"`）
3. **重新登录前端** — 制造执行菜单下出现"工序报工"

> ⚠️ **老工单注意**：5b 之前已下达的工单没有工序快照，完工入库会跳过工序校验（兼容逻辑）。如需补工序，把工单作废重建。

---

## 十二、5b 验证清单

### P0 主流程
- [ ] 建工艺路线：成品 D，3 道工序 OP010 切削 / OP020 钻孔 / OP030 装配，各挂一个工位
- [ ] 建班组 T1，加 2 名组员（admin + 另一个用户）
- [ ] 新建工单（成品 D，计划 10），**不选工艺路线直接下达 → 应被拒绝**
- [ ] 选工艺路线后下达 → 查看工单，"工序与派工"区块出现 3 道工序，状态全为"未派工"
- [ ] OP010 派给 admin 个人；OP020 派给班组 T1；OP030 派给 admin + 工位
- [ ] 进入"工序报工"页，"我的待报工工序"出现 3 道
- [ ] OP010 报工 合格 8 不良 2 → 状态"进行中"（8 < 10 未完工）
- [ ] OP020 报工 合格 9 → **应被拒绝**（上道合格仅 8）
- [ ] OP020 报工 合格 8 → 通过
- [ ] OP030 报工 合格 8 → 通过
- [ ] 领料（按 BOM 领足）后完工入库 10 → **应被拒绝**（末道合格仅 8）
- [ ] 完工入库 8 → 通过；工单 ActualQty = 8

### P1 作废链
- [ ] 作废 OP010 的报工 → **应被拒绝**（OP020 依赖它）
- [ ] 作废完工入库单 → 通过；再作废 OP030 报工 → 通过；再作废 OP020 → 通过；最后作废 OP010 → 通过，工序状态回退到"已派工"
- [ ] 有报工记录的工单点"作废" → 应被拒绝

### P2 权限
- [ ] 用非 admin、非 T1 成员的账号登录 → "我的待报工工序"为空
- [ ] 该账号直接调 `/mfg/report` 对 OP010 报工 → 应被拒绝（已指派给其他员工）

---

## 十三、下一步：5d 批次追溯

5b 完成后，阶段 5 剩下 5d：

- 完工入库时生成**成品批次号**（工单号即批次）
- 领料出库时记录**原料批次**（来自采购入库批次）
- 一张追溯表：成品批次 → 工单 → 工序报工（谁/何时/工位）→ 原料批次 → 供应商
- 正向追溯（这批原料用到了哪些成品）+ 反向追溯（这个成品用了哪批原料）

工作量约为 5b 的 0.7 倍，之后进入阶段 6 报表看板。
