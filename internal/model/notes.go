package model

import "gorm.io/gorm"

// Note 笔记模型
type Note struct {
	gorm.Model        // 继承 ID/CreatedAt/UpdatedAt/DeletedAt
	Title      string `gorm:"type:varchar(100);not null;comment:'笔记标题'"`
	Content    string `gorm:"type:text;not null;comment:'笔记内容'"`
	Category   string `gorm:"type:varchar(50);default:'默认';comment:'笔记分类'"` // 新增分类字段
	UserID     uint   `gorm:"not null;comment:'所属用户ID'"`
	Tags       []Tag  `gorm:"many2many:note_tags;comment:'关联的标签'"` // 多对多（通过中间表 note_tags）
}
