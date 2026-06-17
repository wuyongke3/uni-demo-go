package exam_paper

import "time"

// ExamPaper 试卷
type ExamPaper struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	No        string `json:"no" gorm:"type:varchar(30);uniqueIndex;not null;comment:试卷编号" binding:"required,max=30"` // 试卷编号
	URL       string `json:"url" gorm:"type:varchar(500);comment:试卷链接URL" binding:"omitempty,url,max=500"`            // 试卷链接URL
	Type      string `json:"type" gorm:"type:varchar(20);default:'online';comment:试卷类型" binding:"omitempty,oneof=online offline"` // 试卷类型: online/offline
	FileID    string `json:"file_id" gorm:"type:varchar(100);comment:试卷文件ID" binding:"omitempty,max=100"` // 试卷文件ID
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ExamPaper) TableName() string { return "exam_papers" }
