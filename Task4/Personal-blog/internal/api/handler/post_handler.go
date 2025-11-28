package handler

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/response"
	"Personal-blog/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		postService: service.NewPostRepoService(),
	}
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  uint   `json:"userId" binding:"required"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	post := &model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  req.UserID,
	}
	err := h.postService.CreatePost(post)
	if err != nil {
		response.BadRequest(c, "创建文章失败: "+err.Error())
	}
	response.Success(c, post, "文章创建成功")
}

func (h *PostHandler) GetPost(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	fmt.Sscanf(idParam, "%d", &id)
	var post *model.Post
	post, err := h.postService.GetPostByID(id)
	if err != nil {
		response.BadRequest(c, "获取文章失败: "+err.Error())
		return
	}
	response.Success(c, post, "获取文章成功")
}

type UpdatePostRequest struct {
	Id      uint   `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	var post *model.Post
	post, err := h.postService.GetPostByID(req.Id)
	if err != nil {
		response.BadRequest(c, "更新时获取文章失败"+err.Error())
		return
	}
	value, ok := c.Get("userID")
	//当前登录的用户id跟 文章的用户id不一致 不允许修改
	if !ok || post.UserID != value {
		response.BadRequest(c, "只允许更新自己的文章")
		return
	}
	post.Title = req.Title
	post.Content = req.Content
	err = h.postService.UpdatePost(post)
	if err != nil {
		response.BadRequest(c, "更新文章失败"+err.Error())
		return
	}
	response.Success(c, post, "更新文章成功")
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	fmt.Sscanf(idParam, "%d", &id)
	value, ok := c.Get("userID")
	var post *model.Post
	post, err := h.postService.GetPostByID(id)
	if err != nil {
		response.BadRequest(c, "删除时获取文章失败"+err.Error())
		return
	}
	//当前登录的用户id跟 文章的用户id不一致 不允许修改
	if !ok || post.UserID != value {
		response.BadRequest(c, "只允许删除自己的文章")
		return
	}
	if err := h.postService.DeletePost(id); err != nil {
		response.BadRequest(c, "删除文章失败"+err.Error())
		return
	}
	response.Success(c, id, "删除文章成功")
}

type FindPostReq struct {
	Title string `json:"title"`
}

func (h *PostHandler) FindAll(c *gin.Context) {
	var req FindPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	postModel := &model.Post{
		Title: req.Title,
	}
	post, err := h.postService.GetAllPosts(postModel)
	if err != nil {
		response.BadRequest(c, "获取文章列表失败"+err.Error())
		return
	}
	response.Success(c, post, "获取文章列表成功")
}
