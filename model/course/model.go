package course

import (
	"encoding/json"
	"time"
)

// Course 课程
type Course struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"type:varchar(100);not null;comment:课程名称" binding:"required,max=100"`   // 课程名称
	No             string         `json:"no" gorm:"type:varchar(30);uniqueIndex;not null;comment:课程编号" binding:"required,max=30"` // 课程编号
	MainLecturerIDs json.RawMessage `json:"main_lecturer_ids" gorm:"type:json;comment:主讲师ID列表"`        // 主讲师ID数组 (JSON)
	Category       string         `json:"category" gorm:"type:varchar(50);comment:课程分类" binding:"max=50"`     // 课程分类
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Course) TableName() string { return "courses" }
