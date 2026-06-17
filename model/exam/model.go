package exam

import (
	"time"
	"unigo/model/exam_paper"
	"unigo/model/student"
)

// Exam 考核
type Exam struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Type       string    `json:"type" gorm:"type:varchar(20);not null;comment:考核类型" binding:"required,oneof=schedule course"` // 考核类型: schedule=课表考核 / course=课程考核
	CourseID   *uint     `json:"course_id" gorm:"index;comment:课程ID"`                                                         // 课程ID (课程考核时必填)
	ScheduleID *uint     `json:"schedule_id" gorm:"index;comment:课表ID"`                                                       // 课表ID (课表考核时必填)
	StudentID  uint      `json:"student_id" gorm:"index;not null;comment:学员ID" binding:"required"`                            // 学员ID
	Score      float64   `json:"score" gorm:"type:decimal(5,2);default:0;comment:考核分数" binding:"min=0,max=100"`               // 考核分数
	PaperID    *uint     `json:"paper_id" gorm:"index;comment:试卷ID"`                                                          // 试卷ID
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联字段 (查询时手动填充, 不存库)
	Student *student.Student      `json:"student,omitempty" gorm:"-"`
	Paper   *exam_paper.ExamPaper `json:"paper,omitempty" gorm:"-"`
}

func (Exam) TableName() string { return "exams" }
