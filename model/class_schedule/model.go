package class_schedule

import (
	"encoding/json"
	"time"
	"unigo/model/course"
	"unigo/model/lecturer"
	"unigo/model/student"
)

// ClassSchedule 课表
type ClassSchedule struct {
	ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	No          string          `json:"no" gorm:"type:varchar(60);uniqueIndex;not null;comment:课表编号" binding:"required,max=60"`                     // 课表编号
	CourseID    uint            `json:"course_id" gorm:"index;not null;comment:关联课程ID" binding:"required"`                                          // 关联课程
	StartTime   time.Time       `json:"start_time" gorm:"type:datetime;not null;comment:开课时间" binding:"required"`                                   // 开课时间
	EndTime     time.Time       `json:"end_time" gorm:"type:datetime;not null;comment:结课时间" binding:"required"`                                     // 结课时间
	StudentIDs  json.RawMessage `json:"student_ids" gorm:"type:json;comment:学员ID列表"`                                                                // 学员ID数组 (JSON)
	LecturerIDs json.RawMessage `json:"lecturer_ids" gorm:"type:json;comment:讲师ID列表"`                                                               // 讲师ID数组 (JSON)
	Location    string          `json:"location" gorm:"type:varchar(255);comment:上课地点" binding:"max=255"`                                           // 上课地点(线上链接/线下教室)
	ClassType   string          `json:"class_type" gorm:"type:varchar(20);default:'offline';comment:上课类型" binding:"omitempty,oneof=online offline"` // 上课类型: online/offline
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联字段 (查询时手动填充, 不存库, gorm:"-" 忽略)
	Course    *course.Course      `json:"course,omitempty" gorm:"-"`
	Students  []student.Student   `json:"students,omitempty" gorm:"-"`
	Lecturers []lecturer.Lecturer `json:"lecturers,omitempty" gorm:"-"`
}

func (ClassSchedule) TableName() string { return "class_schedules" }
