package classroom

import "time"

// Classroom 教室
type Classroom struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Building  string `json:"building" gorm:"type:varchar(50);not null;comment:教学楼" binding:"required,max=50"` // 教学楼
	Floor     int    `json:"floor" gorm:"type:int;comment:楼层" binding:"min=0,max=100"`                     // 楼层
	RoomNo    string `json:"room_no" gorm:"type:varchar(20);not null;comment:房间号" binding:"required,max=20"` // 房间号
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Classroom) TableName() string { return "classrooms" }
