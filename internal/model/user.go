package model

import (
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model        // 继承 ID/CreatedAt/UpdatedAt/DeletedAt
	Username   string `gorm:"type:varchar(50);unique;not null;comment:'用户名'"`
	Password   string `gorm:"type:varchar(255);not null;comment:'密码（bcrypt加密）'"`
	Email      string `gorm:"type:varchar(100);unique;not null;comment:'邮箱'"`
	Notes      []Note `gorm:"foreignKey:UserID;references:ID;comment:'关联的笔记'"` // 一对多
}

// BeforeSave 保存前加密密码（钩子函数）
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 密码已加密则跳过（避免更新时重复加密）
	if len(u.Password) == 60 { // bcrypt加密后固定60位
		return nil
	}
	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
