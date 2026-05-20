# GVA 初始状态模块功能说明

> 本文档描述在我们做任何修改之前，gin-vue-admin（以下简称 GVA）原生项目的结构、模块和数据流。
> 这是后续所有修改的"基线"，请务必先看完本篇再看 02、03 等增量文档。

---

## 一、整体技术栈

| 层级 | 技术 |
|---|---|
| 后端语言 | Go 1.22+ |
| Web 框架 | Gin（HTTP 路由 + 中间件） |
| ORM | GORM v2（支持 MySQL/PgSQL/Oracle/MSSQL/SQLite） |
| 鉴权 | JWT（登录态）+ Casbin（接口权限） |
| 缓存 | Redis（可选，多点登录、JWT 黑名单） |
| 配置 | Viper + Fsnotify（读取 `config.yaml`，支持热更新） |
| 日志 | Uber Zap |
| API 文档 | Swagger（swag init 自动生成） |
| 前端框架 | Vue 3 + Pinia + Vue Router 4 |
| 前端 UI | Element Plus |
| 前端构建 | Vite |
| 前端样式 | UnoCSS + SCSS |

---

## 二、后端目录结构（server/）

GVA 后端是**严格分层架构**，自上而下：

```
HTTP 请求
  ↓
router/      路由层：定义 URL，挂载中间件，分发到 api 层
  ↓
middleware/  中间件层：JWT 鉴权、Casbin 权限、操作日志、跨域等
  ↓
api/v1/      接口层（Controller）：参数绑定/校验，调用 service，包装响应
  ↓
service/     服务层（业务逻辑）：核心业务，事务编排，调 model
  ↓
model/       模型层：GORM 数据表结构 + request/response DTO
  ↓
global/      全局对象：DB、Redis、Logger、Config 单例
  ↓
MySQL / Redis
```

### 2.1 关键目录详解

- **`main.go`** — 程序入口，仅一行：`core.RunWindowsServer()` 或 `core.RunServer()`，启动 HTTP。
- **`core/`** — 启动核心：初始化 Viper 配置、Zap 日志、Gin 引擎、HTTP Server。
- **`global/global.go`** — 持有全局单例：`GVA_DB`（GORM 实例）、`GVA_REDIS`、`GVA_LOG`、`GVA_CONFIG`、`GVA_VP`（viper 实例）。
- **`global/model.go`** — 定义 `GVA_MODEL` 基类（`ID/CreatedAt/UpdatedAt/DeletedAt`），所有业务表都嵌入它。
- **`config/`** — `config.yaml` 对应的 Go 结构体。
- **`initialize/`** — 各种"初始化器"：
  - `gorm.go` 根据配置选择 MySQL/PgSQL 等
  - `gorm_biz.go` **业务表自动迁移占位**（这是我们后续要写的入口）
  - `router.go` 注册系统路由
  - `router_biz.go` **业务路由注册占位**（后续要写的入口）
  - `redis.go`、`validator.go`、`timer.go`
- **`middleware/`** — Gin 中间件：
  - `jwt.go` 解析 JWT Token，注入用户信息到 context
  - `casbin_rbac.go` 基于角色的接口级权限
  - `operation.go` 操作日志记录
  - `cors.go`、`limit_ip.go`、`logger.go` 等
- **`model/`** — 业务数据模型：
  - `common/` 通用 DTO：`PageInfo` 分页、`Response` 响应包装、`GetById`
  - `system/` 系统模型：用户、角色、菜单、API、字典、参数、操作日志…
  - `example/` 示例模型：客户、文件上传、附件分类
- **`api/v1/`** — Controller 层：每个 model 对应一个 api 文件，方法接收 `*gin.Context`。
- **`service/`** — 业务层：每个 api 对应一个 service struct，封装 GORM 查询。
- **`router/`** — 路由层：`InitXxxRouter(group)`，挂载到 `PrivateGroup`（需鉴权）或 `PublicGroup`（公开）。
- **`utils/`** — 工具函数：JWT 签发/解析、参数校验（基于反射）、文件上传 OSS 封装、定时器等。
- **`source/`** — 系统初始化数据：超级管理员、默认菜单、默认 API、默认角色（首次启动 `/init/initdb`）。
- **`mcp/`** — 内置 MCP Server，暴露 GVA 操作给 AI 编辑器（Cursor/Cline 等）。
- **`plugin/`** — 插件系统，支持热插拔业务模块。
- **`task/`** — 定时任务。

### 2.2 系统已内置的功能模块（在 `system/` 包中）

| 模块 | 说明 |
|---|---|
| **SysUser** | 用户：登录、密码加密（bcrypt）、绑定多角色 |
| **SysAuthority** | 角色：自上而下的树形结构，支持数据权限范围 |
| **SysBaseMenu / SysMenuBtn** | 菜单 + 按钮级权限，前端动态路由 |
| **SysApi** | 后端 API 注册表，与 Casbin 配合做接口鉴权 |
| **SysDictionary / SysDictionaryDetail** | 字典：枚举值统一管理（如订单状态） |
| **SysParams** | 系统参数：键值对配置 |
| **SysOperationRecord** | 操作日志 |
| **SysLoginLog** | 登录日志 |
| **SysAutoCode / SysAutoCodePackage / SysAutoCodeHistory** | **代码生成器**：根据表结构一键生成 CRUD（前后端） |
| **SysExportTemplate** | Excel 导入导出模板 |
| **JwtBlacklist** | JWT 黑名单（登出后失效） |
| **SysApiToken** | 对外 API Token 签发 |
| **SysVersion** | 版本管理 |
| **SysError** | 错误日志 |

**Example 包** 提供了一组演示业务模块（客户管理、文件上传、断点续传、附件分类），可作为我们写新模块时的参考模板。

---

## 三、前端目录结构（web/src/）

```
main.js              入口，注册 Vue/Pinia/Router/UnoCSS
permission.js        路由守卫：JWT 校验 + 动态路由生成
core/                gva 核心配置（站点配置、欢迎注册）
api/                 后端 API 调用封装（每个后端模块一个 .js）
pinia/               状态管理（user, router, dictionary）
router/              静态路由（登录、初始化、layout 入口）
view/                业务页面
  ├─ dashboard/      首页
  ├─ login/          登录
  ├─ layout/         主框架布局（aside 侧栏、header、setting）
  ├─ superAdmin/     超管：用户/角色/菜单/API/字典/参数/操作日志
  ├─ systemTools/    系统工具：代码生成器、表单生成器、插件、错误日志…
  ├─ example/        示例：客户管理、上传案例
  ├─ person/         个人中心
  └─ init/           首次初始化数据库向导
components/          全局组件
directive/           v-auth 指令（按钮级权限）
utils/
  ├─ request.js      Axios 封装：自动加 Token、统一错误处理
  ├─ asyncRouter.js  动态路由生成
  ├─ btnAuth.js      按钮权限判断
  └─ dictionary.js   字典获取
style/               全局 SCSS
```

### 3.1 前端动态菜单/路由机制（重点）

1. 用户登录 → `/api/user/login` 拿到 JWT
2. 调 `/api/menu/getMenu` 拿到该用户角色绑定的菜单树
3. `permission.js` + `asyncRouter.js` 把菜单树转成 Vue Router 的路由对象
4. 每个菜单项的 `component` 字段对应 `view/` 下的 `.vue` 路径，懒加载
5. 按钮权限：`v-auth="'btnSign'"` 指令读取后端返回的按钮列表，控制显隐

**所以我们后续新增进销存功能，前端只需写 .vue 页面 + 在数据库里挂菜单，不用改静态路由。**

---

## 四、典型请求生命周期（以"创建客户"为例）

1. 前端 `web/src/api/customer.js` 的 `createExaCustomer(data)` 发 `POST /customer/customer`
2. `request.js` 自动在 header 加 `x-token`
3. Gin 命中 `router/example/exa_customer.go` 注册的路由
4. 中间件链：`JWTAuth()`（解析 Token）→ `CasbinHandler()`（接口鉴权）→ `OperationRecord()`（日志）
5. 进入 `api/v1/example/exa_customer.go::CreateExaCustomer`：
   - `c.ShouldBindJSON(&customer)` 绑定参数
   - `utils.Verify(...)` 校验
   - 调 `customerService.CreateExaCustomer(customer)`
6. `service/example/exa_customer.go::CreateExaCustomer` 执行 `global.GVA_DB.Create(&e)`
7. GORM 写 MySQL，返回
8. api 层用 `response.OkWithMessage("创建成功", c)` 返回标准 JSON

**这就是 GVA 的标准 5 层链路，后续我们的所有进销存功能都严格按这个模式来。**

---

## 五、数据库与配置

- 数据库类型在 `config.yaml -> system.db-type` 切换（mysql/pgsql/...）
- 启动时 `initialize/gorm.go::RegisterTables` 会调 `db.AutoMigrate(...)` 自动建表
- 业务表的迁移在 `gorm_biz.go::bizModel()` 里（**目前为空，是我们要填的地方**）
- 业务路由的注册在 `router_biz.go::initBizRouter()` 里（**目前为空，也是我们要填的地方**）
- 首次部署：浏览器访问前端 → 引导到 "初始化数据库" 页面 → 自动建系统表 + 写入超管/默认菜单

---

## 六、依赖注入（GVA 的"全局单例"惯例）

GVA **没有用 DI 框架**，全靠包级变量持有单例：

```go
// 各层的入口对象都是包级 var，在 enter.go 里聚合
var ApiGroupApp = new(ApiGroup)        // api/v1/enter.go
var ServiceGroupApp = new(ServiceGroup) // service/enter.go
var RouterGroupApp = new(RouterGroup)   // router/enter.go

// 子模块在自己的 enter.go 里持有依赖
var customerService = service.ServiceGroupApp.ExampleServiceGroup.CustomerService
```

新增模块时，**必须在 4 个 enter.go 里登记**（model 层无 enter，靠 package 即可）：
1. `api/v1/<mod>/enter.go` 聚合本模块所有 Api
2. `api/v1/enter.go` 把 `<mod>.ApiGroup` 加到顶层 ApiGroup
3. `service/<mod>/enter.go` 同理
4. `service/enter.go` 同理
5. `router/<mod>/enter.go` 同理
6. `router/enter.go` 同理
7. `initialize/router_biz.go` 调用本模块的路由注册
8. `initialize/gorm_biz.go` 把本模块的 model 加入 AutoMigrate

这是 GVA 编码的"八股"，后续我们写每个模块都会遵循。

---

## 七、Go 学习提示（针对本项目重点）

如果你刚学 Go，建议重点看以下语法/概念在 GVA 中是怎么用的：

- **结构体 + 标签（tag）** — `gorm:"comment:..."`、`json:"..."`、`form:"..."`、`binding:"required"`，是 GVA 最常见的代码形态
- **嵌入（embedding）** — 业务结构体 `embed global.GVA_MODEL`，相当于"继承"基类字段
- **方法接收者（receiver）** — `func (s *CustomerService) Xxx()` 是 Go 的"方法"写法
- **包级变量做单例** — 上一节讲的 enter.go 模式
- **错误处理** — `if err != nil { return err }`，Go 没有异常，要逐层返回
- **接口（interface）** — GVA 用得不多，多在 utils 里
- **defer** — 关闭资源、事务回滚常用
- **GORM 链式调用** — `db.Where().Order().Limit().Offset().Find(&list)`
- **Gin 的 Context** — `c.ShouldBindJSON / c.Query / c.Param / c.JSON`

每次我做修改时，会在增量文档里**指出本次涉及到的 Go 语法点**，方便你边做边学。

---

> 下一篇：[02-阶段0-进销存模块设计与骨架.md](./02-阶段0-进销存模块设计与骨架.md)
