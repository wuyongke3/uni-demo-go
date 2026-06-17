# Unigo - 培训机构综合管理平台 (后端)

基于 **Go + Gin + GORM** 的培训机构综合管理系统后端，支持 MySQL / PostgreSQL 双数据库，
采用标准分层架构（Config → Database → Model → Repository → Service → Handler → Router），
提供完整的 RESTful CRUD 接口。

---

## 目录

- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [项目结构详解](#项目结构详解)
- [架构分层说明](#架构分层说明)
- [数据库设计](#数据库设计)
- [API 接口文档](#api-接口文档)
- [配置说明](#配置说明)
- [二次开发指南](#二次开发指南)
  - [新增一个实体模块](#1-新增一个实体模块)
  - [给现有模块添加自定义业务逻辑](#2-给现有模块添加自定义业务逻辑)
  - [添加中间件（鉴权/日志/限流等）](#3-添加中间件鉴权日志限流等)
  - [添加关联查询](#4-添加关联查询)
  - [切换数据库](#5-切换数据库)
- [常见问题](#常见问题)

---

## 环境要求

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | >= 1.22 | 推荐 1.26+ |
| MySQL | >= 8.0 或 MariaDB >= 10.5 | 二选一 |
| PostgreSQL | >= 14 | 二选一 |

---

## 快速开始

### 第一步：克隆项目

```bash
git clone <你的仓库地址>
cd unigo
```

### 第二步：配置数据库

#### 方式 A：使用环境变量（推荐）

**Windows (PowerShell):**
```powershell
$env:DB_DRIVER = "mysql"              # 或 "postgresql"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"                  # PostgreSQL 用 "5432"
$env:DB_USER = "root"                  # PostgreSQL 用 "postgres"
$env:DB_PASSWORD = "你的密码"
$env:DB_NAME = "unigo"
$env:SERVER_PORT = "8000"
```

**Linux / macOS:**
```bash
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=你的密码
export DB_NAME=unigo
export SERVER_PORT=8000
```

#### 方式 B：使用 .env 文件（需配合 godotenv 库）

> 当前版本暂未内置 .env 支持，如需要可在 `config/config.go` 中引入 `github.com/joho/godotenv`。

### 第三步：启动服务

```bash
go run main.go
```

看到以下输出即表示启动成功：
```
2026/01/01 12:00:00 数据库连接成功 (mysql)
2026/01/01 12:00:00 服务器启动于 :8000
```

> 首次启动时 GORM 会**自动建表**，无需手动执行 SQL。

### 第四步：验证接口

```bash
# 查询讲师列表
curl http://localhost:8000/api/v1/lecturers?page=1&page_size=10

# 新增讲师
curl -X POST http://localhost:8000/api/v1/lecturers \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","no":"T001","joined_at":"2024-01-01T08:00:00Z"}'
```

---

## 项目结构详解

```
unigo/
├── main.go                    # ★ 应用入口：加载配置 → 连接DB → 建表 → 注册路由 → 启动HTTP
│
├── config/
│   └── config.go              # 配置管理：从环境变量读取 Server/Database 配置，生成 DSN
│
├── database/
│   └── database.go            # 数据库连接：初始化 GORM 实例（全局单例 database.DB）
│
├── model/                     # ★ 数据模型层：定义每个实体的结构体、GORM 标签、校验标签
│   ├── lecturer.go            #   讲师 (姓名/编号/入职时间/离职时间)
│   ├── student.go             #   学员 (姓名/编号)
│   ├── course.go              #   课程 (名称/编号/主讲师ID数组/分类)
│   ├── class_schedule.go      #   课表 (编号/课程ID/时间/学员列表/讲师列表/地点/类型)
│   ├── classroom.go           #   教室 (教学楼/楼层/房间号)
│   ├── exam.go                #   考核 (类型/课程ID/课表ID/学员ID/分数/试卷ID)
│   ├── exam_paper.go          #   试卷 (编号/URL/类型/文件ID)
│   ├── feedback.go            #   反馈 (留言/学生ID/课程ID/课表ID)
│   └── admin.go               #   管理员 (用户名/姓名/角色)
│
├── repository/                # ★ 数据访问层：封装所有数据库 CRUD 操作
│   └── repository.go          #   baseRepository[T] 泛型基类 + 各实体专用 Repo + AutoMigrate
│
├── service/                   # ★ 业务逻辑层：处理业务规则、调用 Repository
│   └── service.go             #   9 个实体 Service（实现 Operations[T] 接口）
│
├── handler/                   # ★ HTTP 处理层：接收请求 → 校验参数 → 调用 Service → 返回响应
│   └── handler.go             #   CRUDEntity[T] 泛型 Handler（Create/GetByID/List/Update/Delete）
│
├── router/                    # 路由注册层：将 URL 路径绑定到 Handler 方法
│   └── router.go              #   9 个模块 × 5 个 RESTful 接口 = 共 45 个 API
│
├── response/                  # 统一响应工具
│   └── response.go            #   Response 结构体 + Success/Error/PageData 辅助函数
│
├── go.mod                     # Go 模块定义（依赖声明）
├── go.sum                     # 依赖校验和
└── README.md                  # 本文件
```

---

## 架构分层说明

本项目采用经典的 **四层架构 + 配置分离** 模式，数据流如下：

```
HTTP 请求
    │
    ▼
┌─────────────┐
│   Router     │  路由分发：URL → Handler 方法
│  (路由层)     │  文件：router/router.go
└──────┬───────┘
       │
       ▼
┌─────────────┐
│   Handler    │  参数解析/校验 → 调用 Service → 组装响应
│  (处理层)     │  文件：handler/handler.go (泛型 CRUDEntity[T])
└──────┬───────┘
       │
       ▼
┌─────────────┐
│   Service    │  业务逻辑/规则校验 → 调用 Repository
│  (服务层)     │  文件：service/service.go (实现 Operations[T])
└──────┬───────┘
       │
       ▼
┌─────────────┐
│  Repository  │  SQL 操作 (CRUD) ←→ 数据库
│  (仓储层)     │  文件：repository/repository.go (泛型 baseRepository[T])
└──────┬───────┘
       │
       ▼
┌─────────────┐
│   Database   │  GORM 全局实例 (*gorm.DB)
│  (连接层)     │  文件：database/database.go
└─────────────┘
```

**各层职责与规则：**

| 层次 | 职责 | 可以依赖 | 禁止事项 |
|------|------|---------|---------|
| **Router** | URL 路由注册、中间件挂载 | Handler | 不写业务逻辑 |
| **Handler** | 解析请求参数、调用 Service、返回 JSON 响应 | Service, Response | 不直接操作数据库 |
| **Service** | 业务规则、事务控制、数据转换 | Repository, Model | 不涉及 HTTP 相关代码 |
| **Repository** | 封装数据库 CRUD 操作 | Database (gorm.DB), Model | 不包含业务判断 |
| **Model** | 定义数据结构、GORM 标签、校验标签 | 无依赖 | 不含任何方法（纯数据） |
| **Config** | 读取和管理配置项 | 无依赖 | 不硬编码敏感信息 |
| **Response** | 统一 API 响应格式 | 无依赖 | 不做业务处理 |

---

## 数据库设计

### ER 关系概览

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Lecturer │◄───┤  Course  │───►│  Exam    │
│  讲师     │     │   课程    │     │   考核    │
└──────────┘     └────┬─────┘     └────▲─────┘
                      │                │
                      ▼                │
               ┌──────────┐           │
               │ClassSched│───────────┘
               │   课表    │
               └────┬─────┘
                    │
          ┌────────┼────────┐
          ▼        ▼        ▼
    ┌──────────┐ ┌──────┐ ┌─────────┐
    │ Student  │ │Class- │ │Feedback │
    │  学员     │ │ room  │ │  反馈    │
    └──────────┘ └──────┘ └─────────┘
                            │
                            ▼
                       ┌──────────┐
                       │ ExamPaper│
                       │   试卷    │
                       └──────────┘
```

### 表结构一览

| 表名 | 对应模型 | 说明 | 特殊字段 |
|------|---------|------|---------|
| `lecturers` | Lecturer | 讲师 | left_at 可空（离职时间） |
| `students` | Student | 学员 | - |
| `courses` | Course | 课程 | main_lecturer_ids (JSON 数组) |
| `class_schedules` | ClassSchedule | 课表 | student_ids / lecturer_ids (JSON), class_type 枚举 |
| `classrooms` | Classroom | 教室 | building + floor + room_no 组合唯一 |
| `exams` | Exam | 考核 | type 枚举(schedule/course), score decimal(5,2) |
| `exam_papers` | ExamPaper | 试卷 | type 枚举(online/offline) |
| `feedbacks` | Feedback | 反馈 | message text 类型 |
| `admins` | Admin | 管理员 | role 枚举 |

### 字段校验规则

所有模型的字符串字段都通过 `binding:` 标签定义了校验规则：

```go
// 示例：Lecturer 模型的校验标签
Name string `binding:"required,max=50"`      // 必填，最大50字符
No   string `binding:"required,max=30"`      // 必填，最大30字符，uniqueIndex 保证唯一
```

支持的校验规则：

| 标签 | 含义 | 示例 |
|------|------|------|
| `required` | 必填 | `binding:"required"` |
| `max` | 最大长度 | `binding:"max=50"` |
| `min` | 最小值 | `binding:"min=0"` |
| `oneof` | 枚举值 | `binding:"oneof=online offline"` |
| `url` | 合法URL | `binding:"omitempty,url"` |
| `omitempty` | 空值跳过 | 与其他规则组合使用 |

---

## API 接口文档

基础路径：`http://localhost:8000/api/v1`

### 通用响应格式

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**分页列表响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

**错误响应：**
```json
{
  "code": 400,
  "message": "字段 [name] 为必填项"
}
```

### 接口清单

| 模块 | 路径前缀 | POST 创建 | GET 列表 | GET 详情 | PUT 更新 | DELETE 删除 |
|------|---------|----------|---------|---------|---------|------------|
| 讲师 | `/lecturers` | ✅ | ✅?page=&page_size= | ✅/:id | ✅/:id | ✅/:id |
| 学员 | `/students` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 课程 | `/courses` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 课表 | `/class-schedules` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 教室 | `/classrooms` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 考核 | `/exams` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 试卷 | `/exam-papers` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 反馈 | `/feedbacks` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |
| 管理员 | `/admins` | ✅ | ✅ | ✅/:id | ✅/:id | ✅/:id |

### 请求示例

**创建讲师：**
```bash
POST /api/v1/lecturers
Content-Type: application/json

{
  "name": "张三",
  "no": "T001",
  "joined_at": "2024-01-01T08:00:00Z"
}
```

**查询讲师列表（分页）：**
```bash
GET /api/v1/lecturers?page=1&page_size=10
```

**更新讲师：**
```bash
PUT /api/v1/lecturers/1
Content-Type: application/json

{
  "name": "张三丰",
  "no": "T001"
}
```

**删除讲师：**
```bash
DELETE /api/v1/lecturers/1
```

---

## 配置说明

所有配置通过 **环境变量** 注入，在 `config/config.go` 中读取：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `DB_DRIVER` | `mysql` | 数据库类型：`mysql` 或 `postgresql` |
| `DB_HOST` | `localhost` | 数据库主机地址 |
| `DB_PORT` | `3306` | 数据库端口（MySQL 默认 3306，PostgreSQL 默认 5432） |
| `DB_USER` | `root` | 数据库用户名 |
| `DB_PASSWORD` | `` (空) | 数据库密码 |
| `DB_NAME` | `unigo` | 数据库名称 |
| `SERVER_PORT` | `8000` | HTTP 服务监听端口 |

### DSN 生成逻辑

根据 `DB_DRIVER` 自动生成连接字符串：

**MySQL 格式：**
```
user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

**PostgreSQL 格式：**
```
host=host user=user password=password dbname=port port sslmode=disable TimeZone=Asia/Shanghai
```

---

## 二次开发指南

### 1. 新增一个实体模块

以新增「**部门 (Department)**」模块为例，完整步骤如下：

#### Step 1：定义 Model

在 `model/` 下新建 `department.go`：

```go
package model

import "time"

type Department struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string    `json:"name" gorm:"type:varchar(100);not null;comment:部门名称" binding:"required,max=100"`
    Code      string    `json:"code" gorm:"type:varchar(20);uniqueIndex;not null;comment:部门编码" binding:"required,max=20"`
    ParentID  *uint     `json:"parent_id" gorm:"index;comment:上级部门ID"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Department) TableName() string { return "departments" }
```

> **关键点：**
> - `json:"..."` 控制 API 请求/响应的 JSON 字段名
> - `gorm:"..."` 控制数据库列属性（类型、约束、注释）
> - `binding:"..."` 控制入参校验规则
> - 必须实现 `TableName()` 方法指定表名

#### Step 2：定义 Repository

在 `repository/repository.go` 中追加：

```go
// DepartmentRepo 部门仓储
type DepartmentRepo struct{ *baseRepository[model.Department] }

func NewDepartmentRepo() *DepartmentRepo {
    return &DepartmentRepo{NewBaseRepository[model.Department]()}
}
```

同时在 `AutoMigrate()` 函数中注册新模型：

```go
func AutoMigrate() error {
    return database.DB.AutoMigrate(
        // ... 已有模型 ...
        &model.Department{},  // ← 新增这行
    )
}
```

#### Step 3：定义 Service

在 `service/service.go` 中追加：

```go
// --- 部门服务 ---

type DepartmentService struct{ repo *repository.DepartmentRepo }

func NewDepartmentService() *DepartmentService {
    return &DepartmentService{repo: repository.NewDepartmentRepo()}
}

func (s *DepartmentService) Create(entity *model.Department) (*model.Department, error) {
    if err := s.repo.Create(entity); err != nil {
        return nil, errors.New("创建失败: " + err.Error())
    }
    return entity, nil
}

func (s *DepartmentService) GetByID(id uint) (*model.Department, error) {
    e, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("记录不存在") }
        return nil, err
    }
    return e, nil
}

func (s *DepartmentService) List(page, pageSize int) ([]model.Department, int64, error) {
    if page < 1 { page = 1 }
    if pageSize < 1 || pageSize > 100 { pageSize = 10 }
    return s.repo.List(page, pageSize)
}

func (s *DepartmentService) Update(id uint, entity *model.Department) (*model.Department, error) {
    if _, err := s.repo.GetByID(id); err != nil { return nil, err }
    if err := s.repo.Update(id, entity); err != nil { return nil, errors.New("更新失败: " + err.Error()) }
    return entity, nil
}

func (s *DepartmentService) Delete(id uint) error {
    if _, err := s.repo.GetByID(id); err != nil { return err }
    return s.repo.Delete(id)
}
```

#### Step 4：注册路由

在 `router/router.go` 中追加：

```go
// 在 var (...) 块中添加：
departmentHandler = &handler.CRUDEntity[model.Department]{SVC: service.NewDepartmentService()}

// 在 SetupRouter 函数的 api.Group 内部添加：
departments := api.Group("/departments")
{
    departments.POST("", departmentHandler.Create)
    departments.GET("", departmentHandler.List)
    departments.GET("/:id", departmentHandler.GetByID)
    departments.PUT("/:id", departmentHandler.Update)
    departments.DELETE("/:id", departmentHandler.Delete)
}
```

#### Step 5：编译验证

```bash
go build ./...
go run main.go
```

访问 `POST /api/v1/departments` 即可测试新接口。

> **总结：新增一个模块只需改 4 个文件**（model → repository → service → router），Handler 层无需改动（复用泛型 CRUDEntity）。

---

### 2. 给现有模块添加自定义业务逻辑

假设需要给「课表」模块添加一个**按课程ID查询课表列表**的自定义接口：

#### Step 1：在 Repository 添加自定义查询方法

```go
// repository/repository.go 中扩展 ClassScheduleRepo 或新建方法

// FindByCourseID 根据课程ID查询课表列表
func (r *ClassScheduleRepo) FindByCourseID(courseID uint) ([]model.ClassSchedule, error) {
    var schedules []model.ClassSchedule
    err := database.DB.Where("course_id = ?", courseID).Find(&schedules).Error
    return schedules, err
}
```

#### Step 2：在 Service 添加对应方法

```go
// service/service.go 中扩展 ClassScheduleService

func (s *ClassScheduleService) ListByCourse(courseID uint) ([]model.ClassSchedule, error) {
    return s.repo.FindByCourseID(courseID)
}
```

#### Step 3：在 Handler 添加自定义方法

可以在 `handler/handler.go` 中为特定实体扩展，或新建独立 handler 文件：

```go
// handler/class_schedule_handler.go (新建)

package handler

import (
    "net/http"
    "strconv"
    "unigo/response"
    "unigo/service"
    "github.com/gin-gonic/gin"
)

var classScheduleExtHandler struct { svc *service.ClassScheduleService }

// ListByCourse GET /class-schedules/course/:courseId
func (h *classScheduleExtHandler) ListByCourse(c *gin.Context) {
    courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
    if err != nil {
        response.BadRequest(c, "无效的课程ID")
        return
    }
    list, err := h.svc.ListByCourse(uint(courseID))
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, list)
}
```

#### Step 4：注册路由

```go
// router/router.go 中追加：
schedules.GET("/course/:courseId", func(c *gin.Context) {
    classScheduleExtHandler{svc: service.NewClassScheduleService()}.ListByCourse(c)
})
```

---

### 3. 添加中间件（鉴权/日志/限流等）

#### Step 1：创建中间件

在项目根目录或新建 `middleware/` 目录下创建：

```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "unigo/response"
)

// JWTAuth JWT 鉴权中间件示例
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" || !strings.HasPrefix(token, "Bearer ") {
            response.Error(c, http.StatusUnauthorized, 401, "未登录或token已过期")
            c.Abort()
            return
        }

        // TODO: 解析和验证 JWT token
        // claims, err := parseJWT(strings.TrimPrefix(token, "Bearer "))
        // ...

        c.Set("user_id", 1)  // 设置当前用户到上下文
        c.Next()
    }
}
```

#### Step 2：应用到路由

```go
// router/router.go 中修改 SetupRouter：

func SetupRouter(r *gin.Engine) {
    api := r.Group("/api/v1")

    // 公开路由（无需登录）
    public := api.Group("/public")
    {
        // ...
    }

    // 受保护路由（需要鉴权）
    protected := api.Group("")
    protected.Use(middleware.JWTAuth())  // ← 挂载中间件
    {
        // 所有业务路由放在这里
        lecturers := protected.Group("/lecturers")
        // ...
    }
}
```

常用中间件推荐：
- **CORS**: `github.com/gin-contrib/cors`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **限流**: `golang.org/x/time/rate` 或 `github.com/ulule/limiter`
- **请求日志**: Gin 自带 `gin.Logger()` 或 `github.com/gin-contrib/logger`
- **Recovery**: Gin 自带 `gin.Recovery()`（已默认启用）

---

### 4. 添加关联查询

当需要返回关联数据时（如课表中包含课程详情），有两种方式：

#### 方式 A：手动填充（推荐，当前已预留字段）

Model 中已通过 `gorm:"-"` 标记了不存库的关联字段：

```go
// model/class_schedule.go 中已有：
Course    *Course    `json:"course,omitempty" gorm:"-"`
Students  []Student  `json:"students,omitempty" gorm:"-"`
Lecturers []Lecturer `json:"lecturers,omitempty" gorm:"-"`
```

在 Service 中填充：

```go
func (s *ClassScheduleService) GetWithRelations(id uint) (*model.ClassSchedule, error) {
    schedule, err := s.repo.GetByID(id)
    if err != nil { return nil, err }

    // 手动查询关联数据
    course, _ := courseService.GetByID(schedule.CourseID)
    schedule.Course = course

    // 解析 JSON 数组并批量查询学员/讲师
    // ... (根据实际需求实现)

    return schedule, nil
}
```

#### 方式 B：GORM Preload（自动预加载）

如果需要在 Repository 层直接 JOIN 查询：

```go
func (r *ClassScheduleRepo) GetWithPreload(id uint) (*model.ClassSchedule, error) {
    var entity model.ClassSchedule
    err := database.DB.Preload("Course").First(&entity, id).Error
    return &entity, err
}
```

> 注意：方式 B 要求在 Model 中定义 GORM 外键关联（`ForeignKey`/`References` 标签），当前项目未使用此方式以保持灵活性。

---

### 5. 切换数据库

**从 MySQL 切换到 PostgreSQL 只需改一个环境变量：**

```bash
# 原来
set DB_DRIVER=mysql
set DB_PORT=3306

# 改为
set DB_DRIVER=postgresql
set DB_PORT=5432
```

然后重启服务即可。GORM 会自动处理两种数据库的方言差异。

> **注意：** 如果已有数据，迁移前请先备份。不同数据库对 JSON 类型的支持有细微差异（MySQL 用 `JSON` 类型，PostgreSQL 用 `JSONB` 更高效）。

---

## 常见问题

### Q: 启动报错 "database connection failed"

检查以下几点：
1. 数据库服务是否已启动？
2. 环境变量 `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` 是否正确？
3. 数据库 `unigo` 是否已创建？（PostgreSQL 需要预先 `CREATE DATABASE unigo;`，MySQL 会自动创建）

### Q: 如何查看 GORM 生成的 SQL 日志？

`database/database.go` 中已配置 `logger.Info` 级别，启动后终端会输出每条 SQL。如需关闭，将 `logger.Info` 改为 `logger.Silent`。

### Q: 如何添加软删除？

1. 在 Model 中添加 `DeletedAt gorm.DeletedAt \`json:"-" gorm:"index"\``
2. 将 Repository 的 `Delete()` 改为 `database.DB.Delete(new(T), id)`（GORM 会自动转为 UPDATE SET deleted_at=...）
3. 查询时会自动过滤已删除记录

### Q: 分页参数 page/page_size 为空时的行为？

Handler 层设置了默认值 `page=1`, `page_size=10`。Service 层会校正越界值（page<1 归一为 1，page_size 超过 100 限制为 10）。

### Q: 项目能否部署到 Docker？

可以，示例 Dockerfile：

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o unigo main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/unigo .
EXPOSE 8000
CMD ["./unigo"]
```

配合 docker-compose.yml 可一键启动应用 + 数据库。

---

## 技术栈

| 技术 | 用途 | 版本 |
|------|------|------|
| [Gin](https://github.com/gin-gonic/gin) | Web 框架 | v1.12.0 |
| [GORM](https://gorm.io/) | ORM（支持 MySQL / PostgreSQL） | v1.25.x |
| [Validator](https://github.com/go-playground/validator) | 参数校验 | v10.x |
| Go | 语言 | >= 1.22 |

## License

MIT
