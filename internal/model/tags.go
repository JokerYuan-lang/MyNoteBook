package model

import "gorm.io/gorm"

// Tag 标签模型
type Tag struct {
	gorm.Model        // 继承 ID/CreatedAt/UpdatedAt/DeletedAt
	Name       string `gorm:"type:varchar(50);not null;comment:'标签名称'"`
	UserID     uint   `gorm:"not null;comment:'所属用户ID'"`           // 新增用户ID，确保标签按用户隔离
	Notes      []Note `gorm:"many2many:note_tags;comment:'关联的笔记'"` // 多对多
}

// 中间表：笔记-标签关联（无需手动创建，GORM自动生成）

type NoteTag struct {
	NoteID uint `gorm:"primaryKey;comment:'笔记ID'"`
	TagID  uint `gorm:"primaryKey;comment:'标签ID'"`
}
