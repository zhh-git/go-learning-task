package handler

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/response"
	"Personal-blog/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: service.NewCommentService(),
	}
}

type CreateCommentReq struct {
	Content string `json:"content" binding:"required"`
	UserId  uint   `json:"userId" binding:"required"`
	PostId  uint   `json:"postId" binding:"required"`
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	comment := &model.Comment{
		Content: req.Content,
		UserID:  req.UserId,
		PostID:  req.PostId,
	}
	err := h.commentService.CreateComment(comment)
	if err != nil {
		response.BadRequest(c, "创建评论失败: "+err.Error())
		return
	}
	response.Success(c, comment, "评论创建成功")
}

func (h *CommentHandler) FindAllByPostId(c *gin.Context) {
	postId := c.Param("postId")
	var id uint
	fmt.Sscanf(postId, "%d", &id)
	commentList, err := h.commentService.GetCommentsByPostID(id)
	if err != nil {
		response.BadRequest(c, "查询评论列表失败: "+err.Error())
		return
	}
	response.Success(c, commentList, "查询评论列表成功")
}
