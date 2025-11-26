package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JokerYuan-lang/MyNoteBook/internal/model"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NoteService 笔记业务逻辑
type NoteService struct {
	db *gorm.DB
}

// NewNoteService 创建 NoteService 实例
func NewNoteService(db *gorm.DB) *NoteService {
	return &NoteService{db: db}
}

// CreateNote 创建笔记（含标签）
func (s *NoteService) CreateNote(userID uint, title, content, category string, tagNames []string) error {
	// 1. 创建笔记
	note := model.Note{
		Title:    title,
		Content:  content,
		Category: category,
		UserID:   userID,
	}
	if err := s.db.Create(&note).Error; err != nil {
		zap.S().Errorf("创建笔记失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 2. 处理标签（不存在则创建，已存在则关联）
	var tags []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		// 按用户ID+标签名查询（确保标签用户隔离）
		err := s.db.Where("user_id = ? AND name = ?", userID, name).First(&tag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新标签
				tag = model.Tag{Name: name, UserID: userID}
				if err := s.db.Create(&tag).Error; err != nil {
					zap.S().Errorf("创建标签失败: %v", err)
					return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
				}
			} else {
				zap.S().Errorf("查询标签失败: %v", err)
				return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
			}
		}
		tags = append(tags, tag)
	}

	// 3. 关联笔记和标签（多对多）
	if err := s.db.Model(&note).Association("Tags").Replace(&tags); err != nil {
		zap.S().Errorf("关联标签失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return nil
}

// GetNoteList 分页查询笔记列表（支持分类筛选）
func (s *NoteService) GetNoteList(userID uint, page, pageSize int, category string) ([]model.Note, int64, error) {
	var (
		notes []model.Note
		total int64
	)
	category = strings.TrimSpace(category)
	// 构建查询条件（用户ID必选，分类可选）
	db := s.db.Model(&model.Note{}).Where("user_id = ?", userID).Preload("Tags") // Preload 关联查询标签
	if category != "" {
		db = db.Where("category = ?", category)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		zap.S().Errorf("统计笔记总数失败: %v", err)
		return nil, 0, fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 分页查询（offset = (page-1)*pageSize）
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("updated_at DESC").Find(&notes).Error; err != nil {
		zap.S().Errorf("查询笔记列表失败: %v", err)
		return nil, 0, fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return notes, total, nil
}

// GetNoteByID 查询单条笔记
func (s *NoteService) GetNoteByID(userID, noteID uint) (*model.Note, error) {
	var note model.Note
	err := s.db.Where("user_id = ? AND id = ?", userID, noteID).Preload("Tags").First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(errcode.GetMsg(errcode.NotFound))
		}
		zap.S().Errorf("查询笔记失败: %v", err)
		return nil, fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}
	return &note, nil
}

// UpdateNote 更新笔记（含标签）
func (s *NoteService) UpdateNote(userID, noteID uint, title, content, category string, tagNames []string) error {
	// 1. 检查笔记是否存在（且属于当前用户）
	var note model.Note
	err := s.db.Where("user_id = ? AND id = ?", userID, noteID).First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.GetMsg(errcode.NotFound))
		}
		zap.S().Errorf("查询笔记失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 2. 更新笔记基本信息
	note.Title = title
	note.Content = content
	note.Category = category
	if err := s.db.Save(&note).Error; err != nil {
		zap.S().Errorf("更新笔记失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 3. 重新关联标签（先清空旧关联，再关联新标签）
	var tags []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		err := s.db.Where("user_id = ? AND name = ?", userID, name).First(&tag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tag = model.Tag{Name: name, UserID: userID}
				if err := s.db.Create(&tag).Error; err != nil {
					zap.S().Errorf("创建标签失败: %v", err)
					return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
				}
			} else {
				zap.S().Errorf("查询标签失败: %v", err)
				return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
			}
		}
		tags = append(tags, tag)
	}

	// 替换标签关联
	if err := s.db.Model(&note).Association("Tags").Replace(&tags); err != nil {
		zap.S().Errorf("更新标签关联失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return nil
}

// DeleteNote 删除笔记（含关联标签）
func (s *NoteService) DeleteNote(userID, noteID uint) error {
	// 1. 检查笔记是否存在
	var note model.Note
	err := s.db.Where("user_id = ? AND id = ?", userID, noteID).First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(errcode.GetMsg(errcode.NotFound))
		}
		zap.S().Errorf("查询笔记失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 2. 先删除笔记与标签的关联（多对多中间表）
	if err := s.db.Model(&note).Association("Tags").Clear(); err != nil {
		zap.S().Errorf("清空标签关联失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	// 3. 删除笔记
	if err := s.db.Delete(&note).Error; err != nil {
		zap.S().Errorf("删除笔记失败: %v", err)
		return fmt.Errorf(errcode.GetMsg(errcode.ServerError))
	}

	return nil
}
