package feedback

import (
	"time"
	"unigo/model/class_schedule"
	"unigo/model/course"
	"unigo/model/student"
)

// Feedback 反馈
type Feedback struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Message    string    `json:"message" gorm:"type:text;not null;comment:反馈留言" binding:"required"` // 反馈留言
	StudentID  uint      `json:"student_id" gorm:"index;not null;comment:学生ID" binding:"required"`  // 学生ID
	CourseID   *uint     `json:"course_id" gorm:"index;comment:课程ID"`                               // 课程ID
	ScheduleID *uint     `json:"schedule_id" gorm:"index;comment:课表ID"`                             // 课表ID
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联字段 (查询时手动填充, 不存库)
	Student  *student.Student              `json:"student,omitempty" gorm:"-"`
	Course   *course.Course                `json:"course,omitempty" gorm:"-"`
	Schedule *class_schedule.ClassSchedule `json:"schedule,omitempty" gorm:"-"`
}

func (Feedback) TableName() string { return "feedbacks" }
