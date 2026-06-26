package middleware

import (
	"strings"
	"time"

	"openpanda-backend/config"
	"openpanda-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware JWT认证中间件
// 验证请求头中的 Authorization: Bearer <token>
// 验证通过后将用户信息存入 gin.Context，后续 handler 可通过 c.Get("userID") 获取
//
// 使用方式：在 router 中对需要认证的路由组使用
//
//	authGroup := r.Group("/api")
//	authGroup.Use(middleware.JWTAuthMiddleware())
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort() // 阻止后续处理
			return
		}

		// 2. 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. 解析并验证 Token
		cfg := config.LoadConfig()
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法是否为 HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			utils.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 4. 提取 Claims 中的用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Unauthorized(c, "Token解析失败")
			c.Abort()
			return
		}

		// 5. 将用户ID存入上下文，后续 handler 可取出使用
		userID := uint(claims["user_id"].(float64))
		c.Set("userID", userID)

		c.Next() // 继续执行后续 handler
	}
}

// GenerateToken 生成JWT Token（工具函数，供登录接口调用）
// 参数 userID: 用户ID
// 返回: token字符串和错误
func GenerateToken(userID uint) (string, error) {
	cfg := config.LoadConfig()

	// 创建 Claims（载荷）
	now := time.Now() // 当前时间
	claims := jwt.MapClaims{
		"user_id": userID,
		// exp: 过期时间 = 当前时间 + 配置的小时数
		// time.Duration(cfg.JWT.ExpireHour) 将 int 转为 Duration，再乘以 Hour
		"exp": jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHour) * time.Hour)),
		// iat: 签发时间 = 当前时间
		"iat": jwt.NewNumericDate(now),
	}

	// 签名生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GetCurrentUserID 辅助函数：从 gin.Context 中获取当前登录用户ID
// 在 controller 中使用此函数获取用户信息
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}


