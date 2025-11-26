package api

import (
	"strconv"

	"github.com/JokerYuan-lang/MyNoteBook/internal/service"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/errcode"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/response"
	"github.com/JokerYuan-lang/MyNoteBook/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建笔记请求参数

type CreateNoteRequest struct {
	Title    string   `json:"title" binding:"required,max=100"` // 标题最多100位
	Content  string   `json:"content" binding:"required"`       // 内容必填
	Category string   `json:"category" binding:"max=50"`        // 分类最多50位
	TagNames []string `json:"tag_names" binding:"required"`     // 标签必填（至少一个）
}

// 更新笔记请求参数

type UpdateNoteRequest struct {
	NoteID   uint     `form:"note_id" binding:"required,min=1"` // 笔记ID
	Title    string   `form:"title" binding:"required,max=100"` // 标题
	Content  string   `form:"content" binding:"required"`       // 内容
	Category string   `form:"category" binding:"max=50"`        // 分类
	TagNames []string `form:"tag_names" binding:"required"`     // 标签
}

// 笔记列表请求参数（分页+筛选）

type NoteListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`             // 页码（至少1）
	PageSize int    `form:"page_size" binding:"required,min=1,max=50"` // 每页数量（1-50）
	Category string `form:"category,omitempty"`                        // 分类（可选）
}

// NoteAPI 笔记接口
type NoteAPI struct {
	noteService *service.NoteService
}

// NewNoteAPI 创建 NoteAPI 实例
func NewNoteAPI(noteService *service.NoteService) *NoteAPI {
	return &NoteAPI{noteService: noteService}
}

// CreateNote 创建笔记接口
func (a *NoteAPI) CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParam, validator.GetErrorMsg(err))
		return
	}

	// 从上下文获取用户ID（AuthCheck 中间件已存入）
	userID, _ := c.Get("user_id")

	// 调用业务逻辑
	err := a.noteService.CreateNote(userID.(uint), req.Title, req.Content, req.Category, req.TagNames)
	if err != nil {
		response.Error(c, errcode.ServerError, err.Error())
		return
	}

	response.SuccessWithoutData(c)
}

// GetNoteList 分页查询笔记列表接口
func (a *NoteAPI) GetNoteList(c *gin.Context) {
	var req NoteListRequest
	if err := c.ShouldBindQuery(&req); err != nil { // 绑定查询参数（URL参数）
		//response.Error(c, errcode.InvalidParam, validator.GetErrorMsg(err))
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Page参数不合法: " + err.Error(), // 输出具体错误原因
			"data": nil,
		})
		return
	}
	zap.S().Info("page", req.Page)
	userID, _ := c.Get("user_id")
	notes, total, err := a.noteService.GetNoteList(userID.(uint), req.Page, req.PageSize, req.Category)
	if err != nil {
		response.Error(c, errcode.ServerError, err.Error())
		return
	}

	// 返回列表+分页信息
	response.Success(c, gin.H{
		"list":      notes,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetNoteByID 查询单条笔记接口
func (a *NoteAPI) GetNoteByID(c *gin.Context) {
	// 从URL路径获取 note_id
	noteIDStr := c.Query("note_id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParam, "笔记ID格式错误")
		return
	}

	userID, _ := c.Get("user_id")
	note, err := a.noteService.GetNoteByID(userID.(uint), uint(noteID))
	if err != nil {
		if err.Error() == errcode.GetMsg(errcode.NotFound) {
			response.ErrorWithDefaultMsg(c, errcode.NotFound)
		} else {
			response.Error(c, errcode.ServerError, err.Error())
		}
		return
	}

	response.Success(c, note)
}

// UpdateNote 更新笔记接口
func (a *NoteAPI) UpdateNote(c *gin.Context) {
	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParam, validator.GetErrorMsg(err))
		return
	}

	userID, _ := c.Get("user_id")
	err := a.noteService.UpdateNote(userID.(uint), req.NoteID, req.Title, req.Content, req.Category, req.TagNames)
	if err != nil {
		if err.Error() == errcode.GetMsg(errcode.NotFound) {
			response.ErrorWithDefaultMsg(c, errcode.NotFound)
		} else {
			response.Error(c, errcode.ServerError, err.Error())
		}
		return
	}

	response.SuccessWithoutData(c)
}

// DeleteNote 删除笔记接口
func (a *NoteAPI) DeleteNote(c *gin.Context) {
	noteIDStr := c.Query("note_id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		response.Error(c, errcode.InvalidParam, "笔记ID格式错误")
		return
	}

	userID, _ := c.Get("user_id")
	err = a.noteService.DeleteNote(userID.(uint), uint(noteID))
	if err != nil {
		if err.Error() == errcode.GetMsg(errcode.NotFound) {
			response.ErrorWithDefaultMsg(c, errcode.NotFound)
		} else {
			response.Error(c, errcode.ServerError, err.Error())
		}
		return
	}

	response.SuccessWithoutData(c)
}
