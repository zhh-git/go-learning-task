package router

import (
	"Personal-blog/internal/api/handler"
	"Personal-blog/internal/pkg/logger"
	"net/http"
	"strings"
	"time"

	"Personal-blog/internal/pkg/jwt"
	"Personal-blog/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 初始化 Gin 引擎（开发环境启用 debug 模式，生产环境禁用）
	r := gin.Default()

	// 全局中间件：跨域（可选）
	r.Use(Cors())

	// 全局中间件：请求日志（可选）
	r.Use(RequestLogger())

	// 注册 API 路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由
		userHandler := handler.NewUserHandler()
		user := api.Group("/users")
		{
			user.POST("", userHandler.CreateUser)        // POST /api/v1/users 创建用户
			user.GET("/:id", userHandler.GetUserByID)    // GET /api/v1/users/:id 查询用户
			user.PUT("/:id", userHandler.UpdateUser)     // PUT /api/v1/users/:id 更新用户
			user.DELETE("/:id", userHandler.DeleteUser)  // DELETE /api/v1/users/:id 删除用户
			user.POST("/login", userHandler.Login)       // POST /api/v1/users/login 用户登录
			user.POST("/register", userHandler.Register) // POST /api/v1/users/register 用户注册
		}

		//文章相关路由
		postHandler := handler.NewPostHandler()
		post := api.Group("/post")
		{
			post.GET("/:id", postHandler.GetPost)
			post.POST("/findAll", postHandler.FindAll)
			post.Use(JWTAuth()).POST("/add", postHandler.CreatePost) //新增文章
			post.GET("/delete/:id", postHandler.DeletePost)
			post.POST("/update", postHandler.UpdatePost)
		}

		//评论相关路由
		commentHandler := handler.NewCommentHandler()
		comment := api.Group("/comment")
		{
			comment.GET("/:postId", commentHandler.FindAllByPostId)
			comment.Use(JWTAuth()).POST("/add", commentHandler.CreateComment)
		}
	}

	return r
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录请求结束时间和耗时
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 日志内容：方法、路径、状态码、耗时、客户端IP
		logger.Info(
			"method:", c.Request.Method,
			" path:", c.Request.URL.Path,
			" status:", c.Writer.Status(),
			" latency:", latency,
			" ip:", c.ClientIP(),
		)
	}
}

// Auth 中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, "未提供认证 Token")
			c.Abort() // 终止请求流程
			return
		}
		// 2. 校验 Token 格式（必须是 Bearer XXX 格式）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Fail(c, http.StatusUnauthorized, "Token 格式错误（应为 Bearer <token>）")
			c.Abort()
			return
		}
		// 3. 解析并验证 Token
		tokenStr := parts[1]
		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			logger.Warn("Token 验证失败：", err.Error(), "，IP：", c.ClientIP())
			response.Fail(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		// 4. 将用户信息注入上下文（后续接口可通过 c.Get 获取）
		c.Set("jwt_claims", claims)
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 5. 继续执行后续中间件/接口
		c.Next()
	}
}
