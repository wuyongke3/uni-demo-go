package repository

import (
	"fmt"
	"log"
	"unigo/database"
	"unigo/model/admin"
	"unigo/model/class_schedule"
	"unigo/model/classroom"
	"unigo/model/course"
	"unigo/model/exam"
	"unigo/model/exam_paper"
	"unigo/model/feedback"
	"unigo/model/lecturer"
	"unigo/model/student"
)

// ============================================================
//  第一部分：自动建表初始化机制
// ============================================================

// tableRegistry 所有需要自动建表的模型注册表
var tableRegistry = []struct {
	name  string // 表名 (用于日志)
	model any    // GORM 模型实例
}{
	{"lecturers", &lecturer.Lecturer{}},
	{"students", &student.Student{}},
	{"courses", &course.Course{}},
	{"class_schedules", &class_schedule.ClassSchedule{}},
	{"classrooms", &classroom.Classroom{}},
	{"exams", &exam.Exam{}},
	{"exam_papers", &exam_paper.ExamPaper{}},
	{"feedbacks", &feedback.Feedback{}},
	{"admins", &admin.Admin{}},
}

// AutoMigrate 自动建表初始化
//
// 执行流程:
//  1. Ping 数据库，验证连接可用性
//  2. 遍历 tableRegistry，逐张检查表是否存在
//  3. 不存在 → CreateTable 建表（含索引/约束）
//  4. 已存在   → 跳过，不做任何修改
//
// 保证幂等性: 无论执行多少次，结果一致；不会删除或修改已有数据/字段。
func AutoMigrate() error {
	db := database.DB

	// Step 1: 检查数据库连接是否可用 (Ping)
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接不可用 (请检查 DB_HOST/DB_PORT/DB_USER/DB_PASSWORD): %w", err)
	}

	log.Println("[Migration] 数据库连接验证通过，开始检查表结构...")

	var createdCount int
	var skippedCount int

	for _, entry := range tableRegistry {
		migrator := db.Migrator()

		// 检查表是否已存在
		exists := migrator.HasTable(entry.model)

		if exists {
			log.Printf("[Migration] %-20s 已存在，跳过", entry.name)
			skippedCount++
			continue
		}

		// 表不存在，执行建表 (自动创建主键/索引/唯一约束等)
		if createErr := migrator.CreateTable(entry.model); createErr != nil {
			return fmt.Errorf("创建表 [%s] 失败: %w", entry.name, createErr)
		}
		log.Printf("[Migration] %-20s 创建成功 ✓", entry.name)
		createdCount++
	}

	log.Printf("[Migration] 初始化完成: 新建 %d 张表，跳过 %d 张已存在的表",
		createdCount, skippedCount)

	return nil
}

// ============================================================
//  第二部分：泛型 CRUD 仓储层
// ============================================================

// GenericRepository 通用 CRUD 仓储接口 (泛型)
type GenericRepository[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	GetByIDs(ids []uint) ([]T, error)
	List(page, pageSize int) ([]T, int64, error)
	Update(id uint, entity *T) error
	Delete(id uint) error
	BatchDelete(ids []uint) error
}

// baseRepository 通用 CRUD 实现 (基于 GORM)
type baseRepository[T any] struct{}

// NewBaseRepository 创建通用仓储实例
func NewBaseRepository[T any]() *baseRepository[T] {
	return &baseRepository[T]{}
}

// Create 新增记录
func (r *baseRepository[T]) Create(entity *T) error {
	return database.DB.Create(entity).Error
}

// GetByID 根据ID查询
func (r *baseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := database.DB.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 分页查询列表
func (r *baseRepository[T]) List(page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	db := database.DB.Model(new(T))
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}

// Update 更新记录 (只更新非零值字段)
func (r *baseRepository[T]) Update(id uint, entity *T) error {
	return database.DB.Model(new(T)).Where("id = ?", id).Updates(entity).Error
}

// Delete 删除记录 (物理删除)
func (r *baseRepository[T]) Delete(id uint) error {
	return database.DB.Delete(new(T), id).Error
}

// GetByIDs 批量根据ID列表查询
func (r *baseRepository[T]) GetByIDs(ids []uint) ([]T, error) {
	var entities []T
	if err := database.DB.Find(&entities, ids).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// BatchDelete 批量删除 (物理删除)
func (r *baseRepository[T]) BatchDelete(ids []uint) error {
	return database.DB.Delete(new(T), ids).Error
}

// --- 各实体专用 Repository ---

// LecturerRepo 讲师仓储
type LecturerRepo struct {
	*baseRepository[lecturer.Lecturer]
}

func NewLecturerRepo() *LecturerRepo { return &LecturerRepo{NewBaseRepository[lecturer.Lecturer]()} }

// StudentRepo 学员仓储
type StudentRepo struct {
	*baseRepository[student.Student]
}

func NewStudentRepo() *StudentRepo { return &StudentRepo{NewBaseRepository[student.Student]()} }

// CourseRepo 课程仓储
type CourseRepo struct{ *baseRepository[course.Course] }

func NewCourseRepo() *CourseRepo { return &CourseRepo{NewBaseRepository[course.Course]()} }

// ClassScheduleRepo 课表仓储
type ClassScheduleRepo struct {
	*baseRepository[class_schedule.ClassSchedule]
}

func NewClassScheduleRepo() *ClassScheduleRepo {
	return &ClassScheduleRepo{NewBaseRepository[class_schedule.ClassSchedule]()}
}

// ClassroomRepo 教室仓储
type ClassroomRepo struct {
	*baseRepository[classroom.Classroom]
}

func NewClassroomRepo() *ClassroomRepo {
	return &ClassroomRepo{NewBaseRepository[classroom.Classroom]()}
}

// ExamRepo 考核仓储
type ExamRepo struct{ *baseRepository[exam.Exam] }

func NewExamRepo() *ExamRepo { return &ExamRepo{NewBaseRepository[exam.Exam]()} }

// ExamPaperRepo 试卷仓储
type ExamPaperRepo struct {
	*baseRepository[exam_paper.ExamPaper]
}

func NewExamPaperRepo() *ExamPaperRepo {
	return &ExamPaperRepo{NewBaseRepository[exam_paper.ExamPaper]()}
}

// FeedbackRepo 反馈仓储
type FeedbackRepo struct {
	*baseRepository[feedback.Feedback]
}

func NewFeedbackRepo() *FeedbackRepo { return &FeedbackRepo{NewBaseRepository[feedback.Feedback]()} }

// AdminRepo 管理员仓储
type AdminRepo struct{ *baseRepository[admin.Admin] }

func NewAdminRepo() *AdminRepo { return &AdminRepo{NewBaseRepository[admin.Admin]()} }
