package admin

import "time"

// Admin 管理员
type Admin struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" binding:"required,min=3,max=50"` // 用户名
	Name      string `json:"name" gorm:"type:varchar(50);not null;comment:姓名" binding:"required,max=50"`                 // 姓名
	Role      string `json:"role" gorm:"type:varchar(20);default:'admin';comment:角色" binding:"omitempty,oneof=admin super_admin operator"` // 角色
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Admin) TableName() string { return "admins" }
