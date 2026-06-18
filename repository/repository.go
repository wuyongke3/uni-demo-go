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

// AutoMigrate 自动建表/字段同步
//
// 执行流程:
//  1. Ping 数据库，验证连接可用性
//  2. 遍历所有模型，调用 GORM AutoMigrate
//     - 表不存在 → 自动创建（含主键/索引/约束）
//     - 表已存在但缺列 → 自动添加新列（不删除/修改已有字段）
//
// 保证幂等性: 无论执行多少次，结果一致；不会删除已有数据或字段。
func AutoMigrate() error {
	db := database.DB

	// Step 1: 检查数据库连接是否可用
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层连接失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接不可用: %w", err)
	}

	log.Println("[Migration] 数据库连接验证通过，开始同步表结构...")

	// Step 2: 使用 AutoMigrate 自动同步（建表 + 添加缺失列）
	if err := db.AutoMigrate(
		&lecturer.Lecturer{},
		&student.Student{},
		&course.Course{},
		&class_schedule.ClassSchedule{},
		&classroom.Classroom{},
		&exam.Exam{},
		&exam_paper.ExamPaper{},
		&feedback.Feedback{},
		&admin.Admin{},
	); err != nil {
		return fmt.Errorf("表结构同步失败: %w", err)
	}

	log.Println("[Migration] 表结构同步完成 ✓")
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

// FindByNo 根据讲师编号查询 (用于登录)
func (r *LecturerRepo) FindByNo(no string) (*lecturer.Lecturer, error) {
	var entity lecturer.Lecturer
	err := database.DB.Where("no = ?", no).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindByName 根据讲师姓名查询 (用于登录)
func (r *LecturerRepo) FindByName(name string) (*lecturer.Lecturer, error) {
	var entity lecturer.Lecturer
	err := database.DB.Where("name = ?", name).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// StudentRepo 学员仓储
type StudentRepo struct {
	*baseRepository[student.Student]
}

func NewStudentRepo() *StudentRepo { return &StudentRepo{NewBaseRepository[student.Student]()} }

// FindByNo 根据学员编号查询 (用于登录)
func (r *StudentRepo) FindByNo(no string) (*student.Student, error) {
	var entity student.Student
	err := database.DB.Where("no = ?", no).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindByName 根据学员姓名查询 (用于登录)
func (r *StudentRepo) FindByName(name string) (*student.Student, error) {
	var entity student.Student
	err := database.DB.Where("name = ?", name).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

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
