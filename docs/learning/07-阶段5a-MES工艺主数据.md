# 07 - 阶段 5a：MES 工艺主数据（工位 + 班组 + 工艺路线）

## 一、本子阶段定位

阶段 4 的 MES 把生产抽象成一个"工单整体完工"的盒子，里面是黑盒。阶段 5 要把盒子打开——**精确到工序级别管理派工/报工/质检**。

5a 是 5 系列的**主数据基础**，不动业务流程，只补充：
- 工序由谁做？→ **班组 + 工位**
- 工序怎么做？→ **工艺路线**

5b 才会真正改造工单，引入工序级派工/报工。

---

## 二、决策汇总

| 维度 | 选择 | 含义 |
|------|------|------|
| 工序组织 | **工艺路线模式** | 工艺路线作为成品的配套主数据，多版本，工单引用版本快照 |
| 报工 | 员工独立报工（5b 实现） | 每个员工 PDA/网页登录自己报，记录 who/when/howmuch |
| 质检 | **暂不做** | 跳过 5c |
| 追溯 | 批次（5d 实现） | 工单=批次 |
| SPC | **暂不做** | 跳过 5e |
| 节奏 | **5a → 5b → 5d** 分批交付 | 每子阶段验证通过再开下一阶段 |

---

## 三、新增表（5 张）

```
mfg_workstation      工位（人工/设备/外协）
mfg_team             班组
mfg_team_member      班组成员（user ↔ team）
mfg_routing          工艺路线头（多版本）
mfg_routing_item     工艺路线行 = 工序定义
```

### 关键字段

**MfgWorkstation**：
- `Code` 唯一
- `Type` = 1 人工 / 2 设备 / 3 外协
- `Capacity` 标称小时产能（用于 5b APS 雏形）
- `WarehouseID` 可选关联仓库（追踪在哪个车间）

**MfgTeam / MfgTeamMember**：
- 班组扁平结构，`LeaderID` 关联 sys_user.ID
- 多对多关系用独立关联表，避免 GORM 自动命名表
- 成员同样关联 sys_user，**复用 GVA 现有用户体系**

**MfgRouting / MfgRoutingItem**：
- 同一成品多版本（v1.0、v1.1）
- `Seq` 工序号 10/20/30 间隔，方便后期插入新工序
- `Code` OP010 唯一编码，便于追溯
- `StdHours` 标准工时（小时/单位），用于 5b 的工时统计
- `IsKey` 关键工序标记，5d 追溯阶段会重点跟踪

---

## 四、对现有工单的扩展

只加一个字段、一行 import：

```go
// mfg_work_order
RoutingID uint  `json:"routingId" gorm:"comment:引用的工艺路线版本ID（5a扩展，可空）"`
Routing   *MfgRouting `json:"routing,omitempty" gorm:"foreignKey:RoutingID"`
```

**为什么是可选**：
- 已存在的工单没有 RoutingID，强制必填会导致旧数据失效
- 5b 才会强制要求"下达工单时必须有工艺路线"
- 当前**仅作为关联 + 前端下拉展示**，不影响下达/领料/完工逻辑

---

## 五、Service 设计要点

### 5.1 工位删除限制

```go
var cnt int64
global.GVA_DB.Model(&MfgRoutingItem{}).Where("workstation_id=?", id).Count(&cnt)
if cnt > 0 { return errors.New("该工位已被工艺路线引用，不可删除（请改为停用）") }
```

业内通用约束：**主数据被业务引用后只能停用，不能物理删除**——这是审计/历史可追溯的根本要求。

### 5.2 班组的成员是关联表，编辑要全量替换

```go
return tx.Transaction(func(tx *gorm.DB) error {
    tx.Where("team_id=?", e.ID).Delete(&MfgTeamMember{})
    e.Members = nil // 先清空避免 Save 双写
    tx.Save(e)
    for i := range members {
        members[i].ID = 0
        members[i].TeamID = e.ID
        tx.Create(&members[i])
    }
})
```

**为什么不用 FullSaveAssociations**：班组成员是关联表（user_id + team_id），不是子项。FullSave 会试图给每个 member 也写 user 字段，导致 GVA_USER 被错误"更新"。手动 delete + create 最安全。

### 5.3 工艺路线删除：被工单引用后禁删

```go
global.GVA_DB.Model(&MfgWorkOrder{}).Where("routing_id=?", id).Count(&cnt)
if cnt > 0 { return errors.New("该工艺路线已被工单引用，不可删除（请改为停用）") }
```

与 BOM 规则一致：版本快照模式下，源数据要长期保留。

### 5.4 GetByID 的预加载排序

```go
.Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("seq") })
```

**Preload 接受闭包做查询定制**，让 items 永远按 seq 排序，前端不用再手动排。

---

## 六、新增/修改文件清单

### 后端

```
server/
├── model/manufacturing/
│   ├── mfg_master.go                 ★ 新增（工位/班组/工艺路线 5 表）
│   └── mfg_documents.go              ◯ 修改：MfgWorkOrder 加 RoutingID
├── service/manufacturing/
│   ├── master.go                     ★ 新增（3 个 service）
│   └── enter.go                      ◯ 修改：注册新 service
├── api/v1/manufacturing/
│   ├── master.go                     ★ 新增（3 类 API，共 19 端点）
│   └── enter.go                      ◯ 修改
├── router/manufacturing/
│   ├── master.go                     ★ 新增
│   └── enter.go                      ◯ 修改
├── initialize/
│   ├── gorm_biz.go                   ◯ 修改：AutoMigrate 5 表
│   └── router_biz.go                 ◯ 修改：挂 MfgMasterRouter
```

### 前端

```
web/src/
├── api/manufacturing/index.js        ◯ 修改：加 workstationApi/teamApi/routingApi
└── view/manufacturing/
    ├── workstation/index.vue         ★ 新增
    ├── team/index.vue                ★ 新增
    ├── routing/index.vue             ★ 新增
    └── workOrder/index.vue           ◯ 修改：增加"工艺路线"下拉
```

### SQL + 文档

```
db/manufacturing_5a_menu_perm.sql    ★ 新增：3 子菜单 + 19 API + 19 Casbin
docs/learning/07-阶段5a-MES工艺主数据.md  ★ 本文
```

---

## 七、Go / Vue 学习要点

### 7.1 GORM Preload 闭包定制查询

之前写：
```go
.Preload("Items")
```
现在升级：
```go
.Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("seq") })
```
预加载时直接对子查询施加 ORDER BY，避免后端拿到无序数据后再排。**N+1 查询的预防 + 排序一步到位**。

### 7.2 多对多关联的两种写法

**自动多对多（不推荐）**：
```go
type MfgTeam struct {
    Members []SysUser `gorm:"many2many:mfg_team_users;"`
}
```
GORM 会自动建 `mfg_team_users` 表，但你看不到这张表的代码、也不能加额外字段（如加入时间）。

**显式关联表（推荐）**：
```go
type MfgTeamMember struct {
    TeamID uint
    UserID uint
    JoinedAt time.Time  // 可以加额外字段
}
type MfgTeam struct {
    Members []MfgTeamMember `gorm:"foreignKey:TeamID"`
}
```
关联表是一等公民，未来要加"加入日期/角色"等字段直接扩展。

### 7.3 前端"按需领料"的反向继承思路

工单页面从"先选成品 → 加载该成品的 BOM 列表"扩展为"同时加载工艺路线列表"，用 `Promise.all` 并行：

```js
const [bRes, rRes] = await Promise.all([
  bomApi.getActiveByProduct({ id: pid }),
  routingApi.getActiveByProduct({ id: pid })
])
```

**Vue 3 的 reactive 配合 Promise.all** 让多源数据加载几乎零样板。等价的 React 写法要 useState/useEffect 一堆。

---

## 八、对项目结构的影响

| 维度 | 阶段 4 结束 | 5a 结束 |
|------|------------|---------|
| 数据表 | 29 张（inv 22 + mfg 7） | **34 张**（mfg +5） |
| API 端点 | 121 | **+19 = 140** |
| 前端页面 | inv 10 + mfg 4 | **+ 3 主数据 = 17 页** |
| 学习文档 | 01-06 | **+ 07** |

---

## 九、下一步：5b 派工 + 报工

阶段 5b 将做的事（预告）：

1. **工单展开工序**：下达时根据工艺路线展开成 `mfg_work_order_op` 工序快照
2. **派工单**：把工序派给特定员工 / 班组 / 工位
3. **报工单**：员工记录"在工序 X 上做了 Y 数量、用了 Z 小时"
4. **工序状态机**：未派工 → 已派工 → 进行中 → 已完工
5. **进度联动**：所有工序完工 → 工单 ActualQty 才能在完工入库时计入

5b 工作量约为 5a 的 1.5 倍。完成后整个 MES 就是一个**工序级精度的车间执行系统**了。

---

## 十、5a 验证清单

P0：
- [ ] 创建工位：人工/设备/外协 各一个
- [ ] 创建班组：选 1 个班组长 + 选 3 个组员
- [ ] 创建工艺路线（成品 D，3 道工序：OP010 切削 / OP020 钻孔 / OP030 装配）
- [ ] 编辑工艺路线：调整工序顺序
- [ ] 启用/停用工艺路线
- [ ] 工单页面：选成品后能看到工艺路线下拉

P1：
- [ ] 工位被工艺路线引用后不能删除（提示停用）
- [ ] 工艺路线被工单引用后不能删除
- [ ] 同一成品+版本号唯一性校验
- [ ] 班组成员变更（编辑）后正确显示

完成验证后向我反馈，我开始 5b 工作。
