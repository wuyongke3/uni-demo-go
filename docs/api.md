# API 接口文档

基础路径：`http://localhost:8000/api/v1`

---

## 目录

- [统一响应格式](#统一响应格式)
- [鉴权机制](#鉴权机制)
  - [登录接口](#登录接口)
  - [获取当前用户](#获取当前用户)
  - [Token 使用方式](#token-使用方式)
- [5 个标准接口](#5-个标准接口)
- [接口详细说明](#接口详细说明)
  - [POST /{module}/all - 分页列表查询](#1-post-moduleall---分页列表查询)
  - [GET /{module}/info/:ids - 查询详情](#2-get-moduleinfoids---查询详情)
  - [POST /{module}/add - 新增记录](#3-post-moduleadd---新增记录)
  - [POST /{module}/modify/:id - 编辑记录](#4-post-modulemodifyid---编辑记录)
  - [DELETE /{module}/delete/:ids - 批量删除](#5-delete-moduledeleteids---批量删除)
- [所有模块路由汇总](#所有模块路由汇总)
- [各实体请求体字段说明](#各实体请求体字段说明)

---

## 统一响应格式

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
    "total": 100
  }
}
```

**错误响应（业务错误）：**

```json
{
  "code": 60001,
  "message": "账号或密码错误"
}
```

**错误响应（参数校验失败，含字段级详情）：**

```json
{
  "code": 40001,
  "message": "参数校验失败",
  "details": [
    { "field": "Name", "message": "该字段为必填项" },
    { "field": "Password", "message": "最小长度为 6 个字符" }
  ]
}
```

**完整错误码表：**

| 错误码      | 类别         | 说明                   |
| ----------- | ------------ | ---------------------- |
| `0`         | 成功         | 操作成功               |
| `10001`     | 系统错误     | 系统内部异常           |
| `20001`     | 认证错误     | 未登录 (无Token)       |
| `20002`     | 认证错误     | Token 格式错误         |
| `20003`     | 认证错误     | Token 已过期或无效     |
| `30001`     | 数据错误     | 资源不存在             |
| **`40001`** | **参数错误** | **参数校验失败(通用)** |
| **`40002`** | **参数错误** | **必填字段缺失**       |
| **`40003`** | **参数错误** | **字段格式错误**       |
| **`40004`** | **参数错误** | **字段长度超限**       |
| `50001`     | 服务错误     | 数据库操作失败         |
| `50101`     | 服务错误     | 密码处理失败           |
| `50201`     | 服务错误     | Token 生成失败         |
| `60001`     | 业务错误     | 账号或密码错误         |
| `60002`     | 业务错误     | 注册失败               |
| `60003`     | 业务错误     | 编号已被注册           |

> **设计原则：** HTTP 状态码统一返回 `200`，通过 `code` 字段区分成功/失败。前端只需判断 `code === 0` 即可。

---

## 鉴权机制

所有业务接口（除登录外）均需在请求头中携带 **JWT Token** 进行鉴权。

### 登录接口

#### POST /auth/lecturer/register - 讲师注册 (公开)

```bash
POST /api/v1/auth/lecturer/register
Content-Type: application/json

{
  "name": "张老师",    // 姓名
  "password": "123456" // 密码 (6-30位)
}
```

**响应：**

```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1,
    "name": "张老师",
    "no": "T001"
  }
}
```

> **编号规则：** 服务端自动生成，格式 `T` + 8位序号（如 T00000001、T00000042）

#### POST /auth/student/register - 学员注册 (公开)

```bash
POST /api/v1/auth/student/register
Content-Type: application/json

{
  "name": "李同学",    // 姓名
  "password": "123456" // 密码 (6-30位)
}
```

**响应：**

```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1,
    "name": "李同学",
    "no": "S001"
  }
}
```

> **编号规则：** 服务端自动生成，格式 `S` + 8位序号（如 S00000001、S00000042）

#### POST /auth/lecturer/login - 讲师登录 (公开)

```bash
POST /api/v1/auth/lecturer/login
Content-Type: application/json

{
  "name": "张老师",    // 姓名 (登录账号)
  "password": "123456" // 密码
}
```

**响应：**

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 1,
    "username": "张老师",
    "role": "lecturer",
    "user_info": {
      "id": 1,
      "name": "张老师",
      "no": "T001"
    }
  }
}
```

#### POST /auth/student/login - 学员登录

```bash
POST /api/v1/auth/student/login
Content-Type: application/json

{
  "name": "李同学",    // 姓名 (登录账号)
  "password": "123456" // 密码
}
```

**响应：**

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 1,
    "username": "李同学",
    "role": "student",
    "user_info": {
      "id": 1,
      "name": "李同学",
      "no": "S001"
    }
  }
}
```

### 获取当前用户

#### GET /auth/me - 获取当前登录用户信息

```bash
GET /api/v1/auth/me
Authorization: Bearer <token>
```

**响应：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "role": "lecturer",
    "username": "T001"
  }
}
```

### Token 使用方式

所有业务接口（`/{module}/*`）均需要在请求头携带 Token：

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**未携带或 Token 无效时返回：**

```json
{
  "code": 20001,
  "message": "未登录，请先登录"
}
```

---

## 5 个标准接口

每个模块都提供以下 5 个统一接口：

| #     | 方法     | 路径                    | 功能         | 说明                            |
| ----- | -------- | ----------------------- | ------------ | ------------------------------- |
| **1** | `POST`   | `/{module}/all`         | 分页列表查询 | 支持分页、排序、Query 参数      |
| **2** | `GET`    | `/{module}/info/:ids`   | 查询详情     | 支持单 ID 或批量 ID（逗号分隔） |
| **3** | `POST`   | `/{module}/add`         | 新增记录     | 返回新建的完整记录              |
| **4** | `POST`   | `/{module}/modify/:id`  | 编辑记录     | 返回更新后的完整记录            |
| **5** | `DELETE` | `/{module}/delete/:ids` | 批量删除     | 支持单个或批量删除              |

---

## 接口详细说明

### 1. POST /{module}/all - 分页列表查询

**请求方式：** JSON Body 或 Query 参数

```bash
# JSON Body 方式
POST /api/v1/lecturers/all
Content-Type: application/json

{
  "limit": 10,
  "page": 1,
  "start": 0,
  "sort": "-name"
}

# Query 参数方式
POST /api/v1/lecturers/all?limit=10&page=1&start=0&sort=-name
```

**请求参数：**

| 参数    | 类型   | 默认值 | 说明                             |
| ------- | ------ | ------ | -------------------------------- |
| `limit` | int    | 10     | 每页条数 (1-100)                 |
| `page`  | int    | 1      | 页码 (start 优先时忽略)          |
| `start` | int    | 0      | 起始偏移量 (**优先级高于 page**) |
| `sort`  | string | -      | 排序字段 (`-字段名` 表示倒序)    |

**响应示例：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "张三",
        "no": "T001",
        "joined_at": "2024-01-01T00:00:00Z",
        "left_at": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "name": "李四",
        "no": "T002",
        "joined_at": "2024-02-01T00:00:00Z",
        "left_at": null,
        "created_at": "2024-02-01T00:00:00Z",
        "updated_at": "2024-02-01T00:00:00Z"
      }
    ],
    "total": 25
  }
}
```

---

### 2. GET /{module}/info/:ids - 查询详情

**请求方式：** GET

```bash
# 单条查询
GET /api/v1/lecturers/info/1

# 批量查询 (逗号分隔多个 ID)
GET /api/v1/lecturers/info/1,2,3
```

**路径参数：**

| 参数  | 类型   | 说明                        |
| ----- | ------ | --------------------------- |
| `ids` | string | 单个 ID 或逗号分隔的多个 ID |

**响应示例（单条）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "张三",
    "no": "T001",
    "joined_at": "2024-01-01T00:00:00Z",
    "left_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**响应示例（多条）：**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    { "id": 1, "name": "张三", ... },
    { "id": 2, "name": "李四", ... },
    { "id": 3, "name": "王五", ... }
  ]
}
```

---

### 3. POST /{module}/add - 新增记录

**请求方式：** POST JSON

```bash
POST /api/v1/lecturers/add
Content-Type: application/json

{
  "name": "张三",
  "no": "T001",
  "joined_at": "2024-01-01T08:00:00Z"
}
```

**请求体：** 对应实体的 JSON 字段（必填字段见各模型说明）

**响应示例：**

```json
{
  "code": 0,
  "message": "新增成功",
  "data": {
    "id": 1,
    "name": "张三",
    "no": "T001",
    "joined_at": "2024-01-01T08:00:00Z",
    "left_at": null,
    "created_at": "2024-06-17T12:00:00Z",
    "updated_at": "2024-06-17T12:00:00Z"
  }
}
```

**错误响应（参数校验失败）：**

```json
{
  "code": 40001,
  "message": "参数校验失败",
  "details": [{ "field": "Name", "message": "该字段为必填项" }]
}
```

---

### 4. POST /{module}/modify/:id - 编辑记录

**请求方式：** POST JSON

```bash
POST /api/v1/lecturers/modify/1
Content-Type: application/json

{
  "name": "张三丰"
}
```

**路径参数：**

| 参数 | 类型 | 说明            |
| ---- | ---- | --------------- |
| `id` | uint | 要编辑的记录 ID |

**请求体：** 需要更新的字段（**无需传完整数据，只传要改的字段即可**）

**响应示例：**

```json
{
  "code": 0,
  "message": "编辑成功",
  "data": {
    "id": 1,
    "name": "张三丰",
    "no": "T001",
    "joined_at": "2024-01-01T08:00:00Z",
    "left_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-06-17T13:00:00Z"
  }
}
```

---

### 5. DELETE /{module}/delete/:ids - 批量删除

**请求方式：** DELETE

```bash
# 删除单条
DELETE /api/v1/lecturers/delete/1

# 批量删除 (逗号分隔多个 ID)
DELETE /api/v1/lecturers/delete/1,2,3
```

**路径参数：**

| 参数  | 类型   | 说明                        |
| ----- | ------ | --------------------------- |
| `ids` | string | 单个 ID 或逗号分隔的多个 ID |

**响应示例（单条）：**

```json
{
  "code": 0,
  "message": "删除成功"
}
```

**响应示例（批量）：**

```json
{
  "code": 0,
  "message": "批量删除成功 (3 条)"
}
```

---

## 所有模块路由汇总

| 模块   | 路径前缀           | all | info | add | modify | delete |
| ------ | ------------------ | --- | ---- | --- | ------ | ------ |
| 讲师   | `/lecturers`       | ✅  | ✅   | ✅  | ✅     | ✅     |
| 学员   | `/students`        | ✅  | ✅   | ✅  | ✅     | ✅     |
| 课程   | `/courses`         | ✅  | ✅   | ✅  | ✅     | ✅     |
| 课表   | `/class-schedules` | ✅  | ✅   | ✅  | ✅     | ✅     |
| 教室   | `/classrooms`      | ✅  | ✅   | ✅  | ✅     | ✅     |
| 考核   | `/exams`           | ✅  | ✅   | ✅  | ✅     | ✅     |
| 试卷   | `/exam-papers`     | ✅  | ✅   | ✅  | ✅     | ✅     |
| 反馈   | `/feedbacks`       | ✅  | ✅   | ✅  | ✅     | ✅     |
| 管理员 | `/admins`          | ✅  | ✅   | ✅  | ✅     | ✅     |

**总计：9 个模块 × 5 个接口 = 45 个 API**

---

## 各实体请求体字段说明

### 讲师 (Lecturer)

| 字段        | 类型     | 必填 | 说明                    |
| ----------- | -------- | ---- | ----------------------- |
| `name`      | string   | 是   | 姓名 (max 50)           |
| `no`        | string   | 是   | 讲师编号 (max 30, 唯一) |
| `joined_at` | datetime | 否   | 入职时间                |
| `left_at`   | datetime | 否   | 离职时间 (可空)         |

### 学员 (Student)

| 字段   | 类型   | 必填 | 说明                |
| ------ | ------ | ---- | ------------------- |
| `name` | string | 是   | 姓名 (max 50)       |
| `no`   | string | 是   | 编号 (max 30, 唯一) |

### 课程 (Course)

| 字段                | 类型   | 必填 | 说明                    |
| ------------------- | ------ | ---- | ----------------------- |
| `name`              | string | 是   | 课程名称 (max 100)      |
| `no`                | string | 是   | 课程编号 (max 30, 唯一) |
| `main_lecturer_ids` | array  | 否   | 主讲师 ID 列表 (JSON)   |
| `category`          | string | 否   | 分类 (max 50)           |

### 课表 (ClassSchedule)

| 字段           | 类型     | 必填 | 说明                    |
| -------------- | -------- | ---- | ----------------------- |
| `no`           | string   | 是   | 课表编号 (max 60, 唯一) |
| `course_id`    | uint     | 是   | 关联课程 ID             |
| `start_time`   | datetime | 是   | 开课时间                |
| `end_time`     | datetime | 是   | 结课时间                |
| `student_ids`  | array    | 否   | 学员 ID 列表 (JSON)     |
| `lecturer_ids` | array    | 否   | 讲师 ID 列表 (JSON)     |
| `location`     | string   | 否   | 上课地点 (max 255)      |
| `class_type`   | string   | 否   | 类型: online/offline    |

### 教室 (Classroom)

| 字段       | 类型   | 必填 | 说明            |
| ---------- | ------ | ---- | --------------- |
| `building` | string | 是   | 教学楼 (max 50) |
| `floor`    | int    | 否   | 楼层 (0-100)    |
| `room_no`  | string | 是   | 房间号 (max 20) |

### 考核 (Exam)

| 字段          | 类型   | 必填 | 说明                           |
| ------------- | ------ | ---- | ------------------------------ |
| `type`        | string | 是   | 类型: schedule/course          |
| `course_id`   | uint   | 否\* | 课程 ID (type=course 时必填)   |
| `schedule_id` | uint   | 否\* | 课表 ID (type=schedule 时必填) |
| `student_id`  | uint   | 是   | 学员 ID                        |
| `score`       | float  | 否   | 分数 (0-100)                   |
| `paper_id`    | uint   | 否   | 试卷 ID                        |

### 试卷 (ExamPaper)

| 字段      | 类型   | 必填 | 说明                        |
| --------- | ------ | ---- | --------------------------- |
| `no`      | string | 是   | 试卷编号 (max 30, 唯一)     |
| `url`     | string | 否   | 试卷 URL (合法URL, max 500) |
| `type`    | string | 否   | 类型: online/offline        |
| `file_id` | string | 否   | 文件 ID (max 100)           |

### 反馈 (Feedback)

| 字段          | 类型 | 必填 | 说明     |
| ------------- | ---- | ---- | -------- |
| `message`     | text | 是   | 反馈内容 |
| `student_id`  | uint | 是   | 学生 ID  |
| `course_id`   | uint | 否   | 课程 ID  |
| `schedule_id` | uint | 否   | 课表 ID  |

### 管理员 (Admin)

| 字段       | 类型   | 必填 | 说明                             |
| ---------- | ------ | ---- | -------------------------------- |
| `username` | string | 是   | 用户名 (min 3, max 50, 唯一)     |
| `name`     | string | 是   | 姓名 (max 50)                    |
| `role`     | string | 否   | 角色: admin/super_admin/operator |
