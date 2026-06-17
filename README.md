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
- [常见问题](#常见问题)

---

## 环境要求

| 组件       | 版本要求                  | 说明       |
| ---------- | ------------------------- | ---------- |
| Go         | >= 1.22                   | 推荐 1.26+ |
| MySQL      | >= 8.0 或 MariaDB >= 10.5 | 二选一     |
| PostgreSQL | >= 14                     | 二选一     |

---

## 快速开始

### 第一步：克隆项目

```bash
git clone <你的仓库地址>
cd unigo
```

### 第二步：配置数据库

编辑项目根目录下的 `config.yaml` 文件：

```yaml
server:
  port: "8000" # 服务端口

database:
  driver: "mysql" # 数据库类型: mysql / postgresql
  host: "localhost" # 数据库地址
  port: 3306 # 数据库端口 (MySQL:3306, PostgreSQL:5432)
  user: "root" # 用户名
  password: "" # 密码
  dbname: "unigo" # 数据库名
```

> **PostgreSQL 用户注意：** 将 `driver` 改为 `postgresql`，`port` 改为 `5432`。

### 第三步：启动服务

```bash
go run main.go
```

看到以下输出即表示启动成功：

```
========================================
       UniGo 服务启动中...
========================================
  [Server]   端口: 8000
  [Database] 驱动: mysql | 地址: localhost:3306 | 库: unigo
========================================

[Config] 配置加载成功: config.yaml
[Migration] 数据库连接验证通过，开始检查表结构...
[Migration] lecturers              创建成功 ✓
[Migration] students               创建成功 ✓
...

  服务地址: http://localhost:8000
  API 文档: http://localhost:8000/api/v1/lecturers/info/1
========================================
服务启动成功, 监听端口: 8000
```

> **首次启动时 GORM 会自动建表**，无需手动执行 SQL。

### 第四步：验证接口

```bash
# 查询讲师列表 (分页 + 排序)
curl -X POST http://localhost:8000/api/v1/lecturers/all \
  -H "Content-Type: application/json" \
  -d '{"limit": 10, "page": 1, "sort": "-name"}'

# 查询单条详情
curl http://localhost:8000/api/v1/lecturers/info/1

# 批量查询
curl http://localhost:8000/api/v1/lecturers/info/1,2,3

# 新增讲师
curl -X POST http://localhost:8000/api/v1/lecturers/add \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","no":"T001","joined_at":"2024-01-01T08:00:00Z"}'

# 编辑讲师
curl -X POST http://localhost:8000/api/v1/lecturers/modify/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"张三丰"}'

# 删除讲师
curl -X DELETE http://localhost:8000/api/v1/lecturers/delete/1

# 批量删除
curl -X DELETE http://localhost:8000/api/v1/lecturers/delete/1,2,3
```

---

## 项目结构详解

```
unigo/
├── main.go                    # ★ 应用入口：加载配置 → 连接DB → 建表 → 注册路由 → 启动HTTP
│
├── config.yaml                # ★ 配置文件：Server 和 Database 参数 (YAML 格式)
│
├── config/
│   └── config.go              # 配置管理：从 YAML 文件读取配置，生成 DSN
│
├── database/
│   └── database.go            # 数据库连接：初始化 GORM 实例（全局单例 database.DB）
│
├── model/                     # ★ 数据模型层（模块化子包）
│   ├── lecturer/model.go      #   讲师 (姓名/编号/入职时间/离职时间)
│   ├── student/model.go       #   学员 (姓名/编号)
│   ├── course/model.go        #   课程 (名称/编号/主讲师ID数组/分类)
│   ├── class_schedule/model.go # 课表 (编号/课程ID/时间/学员列表/讲师列表/地点/类型)
│   ├── classroom/model.go     #   教室 (教学楼/楼层/房间号)
│   ├── exam/model.go          #   考核 (类型/课程ID/课表ID/学员ID/分数/试卷ID)
│   ├── exam_paper/model.go    #   试卷 (编号/URL/类型/文件ID)
│   ├── feedback/model.go      #   反馈 (留言/学生ID/课程ID/课表ID)
│   └── admin/model.go         #   管理员 (用户名/姓名/角色)
│
├── repository/                # ★ 数据访问层
│   └── repository.go          #   baseRepository[T] 泛型基类 + 各实体专用 Repo + AutoMigrate
│
├── service/                   # ★ 业务逻辑层（模块化子包）
│   ├── service.go             #   CRUDService[T,R] 泛型基类（通用 CRUD 操作）
│   ├── lecturer/service.go    #   讲师服务 (组合泛型基类)
│   ├── student/service.go     #   学员服务
│   ├── course/service.go      #   课程服务
│   ├── class_schedule/service.go # 课表服务
│   ├── classroom/service.go   #   教室服务
│   ├── exam/service.go        #   考核服务
│   ├── exam_paper/service.go  #   试卷服务
│   ├── feedback/service.go    #   反馈服务
│   └── admin/service.go       #   管理员服务
│
├── handler/                   # ★ HTTP 处理层
│   └── handler.go             #   CRUDEntity[T] 泛型 Handler (All/Info/Add/Modify/Delete)
│
├── router/                    # 路由注册层（模块化子包）
│   ├── router.go              #   统一入口：组装各模块路由
│   ├── lecturer/router.go     #   讲师路由
│   ├── student/router.go      #   学员路由
│   ├── course/router.go       #   课程路由
│   ├── class_schedule/router.go # 课表路由
│   ├── classroom/router.go    #   教室路由
│   ├── exam/router.go         #   考核路由
│   ├── exam_paper/router.go   #   试卷路由
│   ├── feedback/router.go     #   反馈路由
│   └── admin/router.go        #   管理员路由
│
├── response/                  # 统一响应工具
│   └── response.go            #   Response 结构体 + Success/Error 辅助函数
│
├── docs/                      # ★ 文档目录
│   └── api.md                 #   API 接口详细文档
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
│  (路由层)     │  文件：router/{module}/router.go
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
│  (服务层)     │  文件：service/service.go (泛型基类) + service/{module}/service.go
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

| 层次           | 职责                                       | 可以依赖                  | 禁止事项               |
| -------------- | ------------------------------------------ | ------------------------- | ---------------------- |
| **Router**     | URL 路由注册、中间件挂载                   | Handler                   | 不写业务逻辑           |
| **Handler**    | 解析请求参数、调用 Service、返回 JSON 响应 | Service, Response         | 不直接操作数据库       |
| **Service**    | 业务规则、事务控制、数据转换               | Repository, Model         | 不涉及 HTTP 相关代码   |
| **Repository** | 封装数据库 CRUD 操作                       | Database (gorm.DB), Model | 不包含业务判断         |
| **Model**      | 定义数据结构、GORM 标签、校验标签          | 无依赖                    | 不含任何方法（纯数据） |
| **Config**     | 读取和管理配置项                           | 无依赖                    | 不硬编码敏感信息       |
| **Response**   | 统一 API 响应格式                          | 无依赖                    | 不做业务处理           |

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

| 表名              | 对应模型      | 说明   | 特殊字段                                           |
| ----------------- | ------------- | ------ | -------------------------------------------------- |
| `lecturers`       | Lecturer      | 讲师   | left_at 可空（离职时间）                           |
| `students`        | Student       | 学员   | -                                                  |
| `courses`         | Course        | 课程   | main_lecturer_ids (JSON 数组)                      |
| `class_schedules` | ClassSchedule | 课表   | student_ids / lecturer_ids (JSON), class_type 枚举 |
| `classrooms`      | Classroom     | 教室   | building + floor + room_no 组合唯一                |
| `exams`           | Exam          | 考核   | type 枚举(schedule/course), score decimal(5,2)     |
| `exam_papers`     | ExamPaper     | 试卷   | type 枚举(online/offline)                          |
| `feedbacks`       | Feedback      | 反馈   | message text 类型                                  |
| `admins`          | Admin         | 管理员 | role 枚举                                          |

### 字段校验规则

所有模型的字符串字段都通过 `binding:` 标签定义了校验规则：

```go
// 示例：Lecturer 模型的校验标签
Name string `binding:"required,max=50"`      // 必填，最大50字符
No   string `binding:"required,max=30"`      // 必填，最大30字符，uniqueIndex 保证唯一
```

支持的校验规则：

| 标签        | 含义     | 示例                             |
| ----------- | -------- | -------------------------------- |
| `required`  | 必填     | `binding:"required"`             |
| `max`       | 最大长度 | `binding:"max=50"`               |
| `min`       | 最小值   | `binding:"min=0"`                |
| `oneof`     | 枚举值   | `binding:"oneof=online offline"` |
| `url`       | 合法URL  | `binding:"omitempty,url"`        |
| `omitempty` | 空值跳过 | 与其他规则组合使用               |

---

## API 接口文档

> 详细接口文档请查看 [docs/api.md](docs/api.md)

**鉴权说明：** 除登录接口外，所有业务接口均需在请求头携带 `Authorization: Bearer <token>`

**公开接口（无需 Token）：**

| 方法   | 路径                      | 功能                    |
| ------ | ------------------------- | ----------------------- |
| `POST` | `/auth/lecturer/register` | 讲师注册                |
| `POST` | `/auth/student/register`  | 学员注册                |
| `POST` | `/auth/lecturer/login`    | 讲师登录                |
| `POST` | `/auth/student/login`     | 学员登录                |
| `GET`  | `/auth/me`                | 获取当前用户 (需 Token) |

**业务接口（需 Token 鉴权）：**

| #     | 方法     | 路径                    | 功能         |
| ----- | -------- | ----------------------- | ------------ |
| **1** | `POST`   | `/{module}/all`         | 分页列表查询 |
| **2** | `GET`    | `/{module}/info/:ids`   | 查询详情     |
| **3** | `POST`   | `/{module}/add`         | 新增记录     |
| **4** | `POST`   | `/{module}/modify/:id`  | 编辑记录     |
| **5** | `DELETE` | `/{module}/delete/:ids` | 批量删除     |

**所有模块路由：**

| 模块   | 路径前缀           |
| ------ | ------------------ |
| 讲师   | `/lecturers`       |
| 学员   | `/students`        |
| 课程   | `/courses`         |
| 课表   | `/class-schedules` |
| 教室   | `/classrooms`      |
| 考核   | `/exams`           |
| 试卷   | `/exam-papers`     |
| 反馈   | `/feedbacks`       |
| 管理员 | `/admins`          |

**总计：9 个模块 × 5 个接口 = 45 个 API**

---

## 配置说明

所有配置通过 **YAML 文件** 管理，位于项目根目录的 `config.yaml`。

### 配置文件模板

```yaml
server:
  port: "8000" # 服务端口

database:
  driver: "mysql" # 数据库类型: mysql / postgresql
  host: "localhost" # 数据库地址
  port: 3306 # 数据库端口 (MySQL:3306, PostgreSQL:5432)
  user: "root" # 用户名
  password: "" # 密码
  dbname: "unigo" # 数据库名

jwt:
  secret: "your-secret-key" # JWT 签名密钥 (生产环境请修改为复杂随机字符串)
  expire_hour: 24 # Token 过期时间 (小时)
```

### 配置项说明

| 配置项              | 类型   | 默认值        | 说明                                                |
| ------------------- | ------ | ------------- | --------------------------------------------------- |
| `server.port`       | string | `"8000"`      | HTTP 服务监听端口                                   |
| `database.driver`   | string | `"mysql"`     | 数据库类型：`mysql` 或 `postgresql`                 |
| `database.host`     | string | `"localhost"` | 数据库主机地址                                      |
| `database.port`     | int    | `3306`        | 数据库端口（MySQL 默认 3306，PostgreSQL 默认 5432） |
| `database.user`     | string | `"root"`      | 数据库用户名                                        |
| `database.password` | string | `` (空)       | 数据库密码                                          |
| `database.dbname`   | string | `"unigo"`     | 数据库名称                                          |
| `jwt.secret`        | string | -             | JWT 签名密钥（生产环境必须修改）                    |
| `jwt.expire_hour`   | int    | `24`          | Token 有效期（小时）                                |

### 配置文件查找顺序

程序启动时会按以下顺序查找配置文件：

1. **环境变量 `CONFIG_PATH` 指定的路径**（仅用于定位文件位置，不覆盖配置值）
2. **当前目录** 下的 `config.yaml` 或 `config.yml`
3. **可执行文件同目录** 下的 `config.yaml` 或 `config.yml`

如果找不到配置文件，将使用内置默认值启动。

### 指定自定义配置文件路径

```bash
# Linux / macOS
CONFIG_PATH=/etc/unigo/config.yaml ./unigo

# Windows PowerShell
$env:CONFIG_PATH = "C:\etc\unigo\config.yaml"; .\unigo.exe
```

### DSN 生成逻辑

根据 `database.driver` 自动生成连接字符串：

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

在 `model/` 下新建 `department/model.go`：

```go
package department

import "time"

// Department 部门
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
>
> - 包名为 `package department`（小写单数）
> - `json:"..."` 控制 API 请求/响应的 JSON 字段名
> - `gorm:"..."` 控制数据库列属性（类型、约束、注释）
> - `binding:"..."` 控制入参校验规则
> - 必须实现 `TableName()` 方法指定表名

#### Step 2：扩展 Repository

在 `repository/repository.go` 中追加：

```go
import (
    // ... 已有 import ...
    "unigo/model/department"  // ← 新增导入
)

// 在 tableRegistry 中追加：
{"departments", &department.Department{}},

// 在文件末尾追加专用 Repo：
// DepartmentRepo 部门仓储
type DepartmentRepo struct{ *baseRepository[department.Department] }

func NewDepartmentRepo() *DepartmentRepo {
    return &DepartmentRepo{NewBaseRepository[department.Department]()}
}
```

#### Step 3：创建 Service 子包

在 `service/` 下新建 `department/service.go`：

```go
package department

import (
    "unigo/model/department"
    "unigo/repository"
    "unigo/service"
)

// Service 部门服务 (组合泛型基类)
type Service struct{ *service.CRUDService[department.Department, *repository.DepartmentRepo] }

// New 创建部门服务实例
func New() *Service {
    return &Service{CRUDService: service.NewCRUDService[department.Department, *repository.DepartmentRepo](repository.NewDepartmentRepo())}
}
```

#### Step 4：创建 Router 子包

在 `router/` 下新建 `department/router.go`：

```go
package department

import (
    "unigo/handler"
    "unigo/model/department"
    svc "unigo/service/department"

    "github.com/gin-gonic/gin"
)

var h = &handler.CRUDEntity[department.Department]{SVC: svc.New()}

// Register 注册部门模块路由
func Register(rg *gin.RouterGroup) {
    rg.POST("all", h.All)
    rg.GET("info/:ids", h.Info)
    rg.POST("add", h.Add)
    rg.POST("modify/:id", h.Modify)
    rg.DELETE("delete/:ids", h.Delete)
}
```

#### Step 5：注册到主路由

在 `router/router.go` 中追加：

```go
import (
    // ... 已有 import ...
    "unigo/router/department"  // ← 新增导入
)

// 在 SetupRouter 函数中追加：
department.Register(api.Group("/departments"))
```

#### Step 6：编译验证

```bash
go build ./...
go run main.go
```

访问 `POST /api/v1/departments/add` 即可测试新接口。

> **总结：新增一个模块只需改 5 个文件**（model → repository → service → router → router 入口），Handler 层无需改动（复用泛型 CRUDEntity）。

---

### 2. 给现有模块添加自定义业务逻辑

假设需要给「课表」模块添加一个**按课程ID查询课表列表**的自定义接口：

#### Step 1：在 Repository 添加自定义查询方法

在 `repository/repository.go` 中扩展 `ClassScheduleRepo` 或新建独立方法：

```go
// FindByCourseID 根据课程ID查询课表列表
func (r *ClassScheduleRepo) FindByCourseID(courseID uint) ([]class_schedule.ClassSchedule, error) {
    var schedules []class_schedule.ClassSchedule
    err := database.DB.Where("course_id = ?", courseID).Find(&schedules).Error
    return schedules, err
}
```

#### Step 2：在 Service 添加对应方法

在 `service/class_schedule/service.go` 中追加：

```go
// ListByCourse 按课程ID查询课表
func (s *Service) ListByCourse(courseID uint) ([]class_schedule.ClassSchedule, error) {
    return s.repo.FindByCourseID(courseID)
}
```

#### Step 3：在 Handler 添加自定义方法

可以在 `handler/handler.go` 中扩展，或新建独立 handler 文件：

```go
// handler/class_schedule_ext.go (新建)

package handler

import (
    "net/http"
    "strconv"
    "unigo/response"

    "github.com/gin-gonic/gin"
)

// ClassScheduleExtHandler 课表扩展处理器
type ClassScheduleExtHandler struct {
    SVC interface {
        ListByCourse(courseID uint) ([]interface{}, error)
    }
}

// ListByCourse GET /class-schedules/course/:courseId
func (h *ClassScheduleExtHandler) ListByCourse(c *gin.Context) {
    courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 32)
    if err != nil {
        response.BadRequest(c, "无效的课程ID")
        return
    }
    list, err := h.SVC.ListByCourse(uint(courseID))
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, list)
}
```

#### Step 4：注册路由

在 `router/class_schedule/router.go` 中追加：

```go
rg.GET("course/:courseId", classScheduleExtHandler.ListByCourse)
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
    "unigo/response"

    "github.com/gin-gonic/gin"
)

// JWTAuth JWT 鉴权中间件示例
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" || !strings.HasPrefix(token, "Bearer ") {
            c.JSON(http.StatusUnauthorized, response.Response{
                Code:    401,
                Message: "未登录或token已过期",
            })
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

修改 `router/router.go` 中的 `SetupRouter` 函数：

```go
func SetupRouter(r *gin.Engine) {
    api := r.Group("/api/v1")

    // 公开路由（无需登录）
    public := api.Group("")
    {
        // 公开的健康检查等接口...
    }

    // 受保护路由（需要鉴权）
    protected := api.Group("")
    protected.Use(middleware.JWTAuth())  // ← 挂载中间件
    {
        // 所有业务路由放在这里
        lecturer.Register(protected.Group("/lecturers"))
        student.Register(protected.Group("/students"))
        // ... 其他模块
    }
}
```

常用中间件推荐：

- **CORS**: `github.com/gin-contrib/cors`
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **限流**: `golang.org/x/time/rate` 或 `github.com/ulule/limiter`
- **请求日志**: Gin 自带 `gin.Logger()` 或 `github.com/gin-contrib/logger`

---

### 4. 添加关联查询

当前项目中 `ClassSchedule`、`Exam`、`Feedback` 等模型包含关联字段（标记为 `gorm:"-"`），需要在查询后手动填充。

以「课表详情带关联数据」为例：

#### Step 1：在 Repository 添加关联查询

```go
// repository/repository.go

// GetWithRelations 查询课表并填充关联数据
func (r *ClassScheduleRepo) GetWithRelations(id uint) (*class_schedule.ClassSchedule, error) {
    var entity class_schedule.ClassSchedule
    if err := database.DB.First(&entity, id).Error; err != nil {
        return nil, err
    }

    // 手动填充关联数据
    database.DB.First(&entity.Course, entity.CourseID)

    // 解析 JSON 数组并查询
    if len(entity.StudentIDs) > 0 {
        json.Unmarshal(entity.StudentIDs, &entity.Students)
    }
    if len(entity.LecturerIDs) > 0 {
        json.Unmarshal(entity.LecturerIDs, &entity.Lecturers)
    }

    return &entity, nil
}
```

#### Step 2：在 Service 暴露方法

```go
// service/class_schedule/service.go

func (s *Service) GetDetail(id uint) (*class_schedule.ClassSchedule, error) {
    return s.repo.GetWithRelations(id)
}
```

#### Step 3：在 Handler 处理请求

```go
// handler/class_schedule_detail.go

func (h *ClassScheduleExtHandler) Detail(c *gin.Context) {
    id, _ := parseUintParam(c, "id")
    data, err := h.SVC.GetDetail(id)
    if err != nil {
        response.NotFound(c, err.Error())
        return
    }
    response.Success(c, data)
}
```

#### Step 4：注册路由

```go
rg.GET("detail/:id", extHandler.Detail)  // GET /class-schedules/detail/1
```

---

### 5. 切换数据库

只需修改 `config.yaml` 中的配置：

```yaml
# 从 MySQL 切换到 PostgreSQL
database:
  driver: "postgresql" # 改为 postgresql
  host: "localhost"
  port: 5432 # 改为 PostgreSQL 默认端口
  user: "postgres" # 改为 PostgreSQL 默认用户
  password: "your_password"
  dbname: "unigo"
```

重启服务即可，GORM 会根据驱动自动选择对应的方言。

---

## 常见问题

### Q: 首次启动报错 "Unknown database 'unigo'"

**A:** 需要先手动创建数据库：

```sql
-- MySQL
CREATE DATABASE unigo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- PostgreSQL
CREATE DATABASE unigo;
```

然后重新启动服务，GORM 会自动建表。

### Q: 如何查看 GORM 生成的 SQL？

**A:** 当前已开启 Info 级别日志，控制台会输出 SQL。如需更详细的日志，可修改 `database/database.go`：

```go
logger.Default.LogMode(logger.Info)  // 改为 logger.Default.LogMode(logger.LogMode(4))
```

### Q: 如何添加软删除功能？

**A:** 在 Model 中添加 `DeletedAt gorm.DeletedAt` 字段即可启用 GORM 的全局软删除：

```go
type Lecturer struct {
    ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string         `json:"name" ...`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`  // ← 添加这行
    // ...
}
```

### Q: 排序字段不存在怎么办？

**A:** 当前实现中，如果排序字段不存在于模型中，排序会被忽略（保持原顺序）。不会报错。

### Q: 如何部署到生产环境？

**A:** 推荐步骤：

1. 交叉编译目标平台二进制：
   ```bash
   GOOS=linux GOARCH=amd64 go build -o unigo-server .
   ```
2. 将 `config.yaml` 和 `unigo-server` 上传到服务器
3. 修改 `config.yaml` 中的数据库连接信息
4. 使用 systemd 或 supervisor 管理进程
5. 配置 Nginx 反向代理（可选）

---

## License

MIT
