package student

import "time"

// Student 学员
type Student struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null;comment:学员姓名" binding:"required,max=50"` // 学员姓名
	No        string    `json:"no" gorm:"type:varchar(30);uniqueIndex;not null;comment:学员编号" binding:"required,max=30"` // 学员编号
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Student) TableName() string { return "students" }
