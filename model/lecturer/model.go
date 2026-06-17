package lecturer

import "time"

// Lecturer 讲师
type Lecturer struct {
	ID        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Password  string     `json:"-" gorm:"type:varchar(100);not null;comment:密码(加密存储)" binding:"required,min=6"`          // 密码 (不返回前端)
	Name      string     `json:"name" gorm:"type:varchar(50);not null;comment:姓名" binding:"required,max=50"`             // 姓名
	No        string     `json:"no" gorm:"type:varchar(30);uniqueIndex;not null;comment:讲师编号" binding:"required,max=30"` // 讲师编号
	JoinedAt  time.Time  `json:"joined_at" gorm:"type:datetime;comment:入职时间"`                                            // 入职时间
	LeftAt    *time.Time `json:"left_at" gorm:"type:datetime;comment:离职时间"`                                              // 离职时间 (可空)
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Lecturer) TableName() string { return "lecturers" }
