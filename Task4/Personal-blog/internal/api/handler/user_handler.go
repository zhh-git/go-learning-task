package handler

import (
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/jwt"
	"Personal-blog/internal/pkg/response"
	"Personal-blog/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.userService.CreateUser(user)
	if err != nil {
		response.BadRequest(c, "创建用户失败: "+err.Error())
		return
	}
	response.Success(c, user, "创建用户成功")
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.BadRequest(c, "获取用户失败: "+err.Error())
		return
	}
	response.Success(c, user, "获取用户成功")
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userService.GetUserByUsername(username)
	if err != nil {
		response.BadRequest(c, "获取用户失败: "+err.Error())
		return
	}
	response.Success(c, user, "获取用户成功")
}

type UpdateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.BadRequest(c, "获取用户失败: "+err.Error())
		return
	}
	user.Email = req.Email
	user.Password = req.Password
	err = h.userService.UpdateUser(user)
	if err != nil {
		response.BadRequest(c, "更新用户失败: "+err.Error())
		return
	}
	response.Success(c, user, "更新用户成功")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	err = h.userService.DeleteUser(id)
	if err != nil {
		response.BadRequest(c, "删除用户失败: "+err.Error())
		return
	}
	response.Success(c, nil, "删除用户成功")
}

func (h *UserHandler) Register(c *gin.Context) {
	var user CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.BadRequest(c, "密码加密失败: "+err.Error())
		return
	}
	user.Password = string(hashedPassword)
	userModel := &model.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	err = h.userService.CreateUser(userModel)
	if err != nil {
		response.BadRequest(c, "注册失败: "+err.Error())
		return
	}
	response.Success(c, user, "注册成功")
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var user LoginUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	var storedUser *model.User
	storedUser, err := h.userService.GetUserByUsername(user.Username)
	if err != nil {
		response.BadRequest(c, "无效的用户名或密码")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		response.BadRequest(c, "无效的用户名或密码")
		return
	}

	// 生成 JWT
	tokenString, err := jwt.GenerateToken(storedUser)
	if err != nil {
		response.BadRequest(c, "生成Token失败")
		return
	}
	// 剩下的逻辑... token放入缓存 等等操作
	response.Success(c, tokenString, "登录成功")
}
