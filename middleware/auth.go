package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ContextKey 当前用户上下文 key
const ContextUserID = "user_id"    // 用户 ID
const ContextRole = "role"         // 角色 (lecturer/student)
const ContextUserName = "username" // 登录账号

// JWTConfigInterface JWT 配置接口 (方便依赖注入)
type JWTConfigInterface interface {
	GetSecret() string
	GetExpireHour() int
}

// JWTAuth JWT 鉴权中间件
//
// 从请求头 Authorization: Bearer <token> 中提取并验证 Token
// 验证通过后将 user_id / role / username 写入 gin.Context
func JWTAuth(cfg JWTConfigInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未登录，请先登录",
			})
			c.Abort()
			return
		}

		// 格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误，格式: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析并验证 Token
		claims, err := parseToken(tokenString, cfg.GetSecret())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token 已过期或无效，请重新登录",
			})
			c.Abort()
			return
		}

		// 将用户信息写入上下文，供后续 Handler 使用
		c.Set(ContextUserID, claims["user_id"])
		c.Set(ContextRole, claims["role"])
		c.Set(ContextUserName, claims["username"])

		c.Next()
	}
}

// ============================================================
//  JWT 工具函数
// ============================================================

// CustomClaims 自定义 JWT 载荷
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"` // lecturer / student
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, role, username string, secret string, expireHour int) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Role:     role,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireHour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "unigo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// parseToken 解析并验证 Token，返回载荷数据
func parseToken(tokenString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
