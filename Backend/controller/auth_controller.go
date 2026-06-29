// ============================================================
// 包名: controller
// 说明: 认证控制器 - 登录、获取当前用户
// ============================================================
package controller

import (
	"openpanda-backend/config"
	"openpanda-backend/middleware"
	"openpanda-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct{}

// NewAuthController 构造函数
func NewAuthController() *AuthController {
	return &AuthController{}
}

// ---------- 请求/响应结构体 ----------

// LoginRequest 登录请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名（必填）
	Password string `json:"password" binding:"required"` // 密码（必填）
}

// LoginResponse 登录响应体
type LoginResponse struct {
	Token    string `json:"token"`    // JWT Token
	Username string `json:"username"` // 用户名
}

// ---------- 接口 ----------

// Login 登录接口
// POST /api/v1/login
// 请求体: { "username": "admin", "password": "xxx" }
// 响应:   { "code": 200, "data": { "token": "...", "username": "admin" } }
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入用户名和密码")
		return
	}

	// 校验管理员凭据（从环境变量读取，默认 admin/!Wo3158023）
	adminUser := config.GetEnv("ADMIN_USERNAME", "admin")
	adminPass := config.GetEnv("ADMIN_PASSWORD", "!Wo3158023")

	if req.Username != adminUser || req.Password != adminPass {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 JWT Token（管理员 userID 固定为 1）
	token, err := middleware.GenerateToken(1)
	if err != nil {
		utils.InternalError(c, "Token 生成失败")
		return
	}

	utils.SuccessWithMessage(c, "登录成功", LoginResponse{
		Token:    token,
		Username: req.Username,
	})
}

// GetMe 获取当前登录用户信息（需 JWT 认证）
// GET /api/v1/auth/me
// 用于前端验证 Token 是否有效，以及获取当前用户信息
func (ctrl *AuthController) GetMe(c *gin.Context) {
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		utils.Unauthorized(c, "未登录")
		return
	}

	utils.Success(c, gin.H{
		"user_id":  userID,
		"username": "admin",
		"role":     "admin",
	})
}
