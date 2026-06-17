package handler

import (
	"net/http"

	"unigo/middleware"
	"unigo/model/lecturer"
	"unigo/model/student"
	"unigo/repository"
	"unigo/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求体
type LoginRequest struct {
	No       string `json:"no" binding:"required"`             // 编号 (讲师编号/学员编号)
	Password string `json:"password" binding:"required,min=6"` // 密码
}

// RegisterRequest 注册请求体 (通用)
type RegisterRequest struct {
	No       string `json:"no" binding:"required,max=30"`             // 编号
	Name     string `json:"name" binding:"required,max=50"`           // 姓名
	Password string `json:"password" binding:"required,min=6,max=30"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`               // JWT Token
	UserID   uint        `json:"user_id"`             // 用户 ID
	Username string      `json:"username"`            // 登录账号
	Role     string      `json:"role"`                // 角色: lecturer / student
	UserInfo interface{} `json:"user_info,omitempty"` // 用户基本信息 (不含密码)
}

// AuthHandler 认证处理器 (持有 JWT 配置)
type AuthHandler struct {
	JWTConfig middleware.JWTConfigInterface
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(cfg middleware.JWTConfigInterface) *AuthHandler {
	return &AuthHandler{JWTConfig: cfg}
}

// LecturerLogin POST /auth/lecturer/login - 讲师登录
func (h *AuthHandler) LecturerLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, formatValidationError(err))
		return
	}

	repo := repository.NewLecturerRepo()
	user, err := repo.FindByNo(req.No)
	if err != nil || user == nil {
		response.BadRequest(c, "账号或密码错误")
		return
	}

	if !checkPassword(req.Password, user.Password) {
		response.BadRequest(c, "账号或密码错误")
		return
	}

	token, err := middleware.GenerateToken(
		user.ID, "lecturer", user.No,
		h.JWTConfig.GetSecret(), h.JWTConfig.GetExpireHour(),
	)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "登录成功",
		Data: LoginResponse{
			Token:    token,
			UserID:   user.ID,
			Username: user.No,
			Role:     "lecturer",
			UserInfo: lecturerInfo(user),
		},
	})
}

// StudentLogin POST /auth/student/login - 学员登录
func (h *AuthHandler) StudentLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, formatValidationError(err))
		return
	}

	repo := repository.NewStudentRepo()
	user, err := repo.FindByNo(req.No)
	if err != nil || user == nil {
		response.BadRequest(c, "账号或密码错误")
		return
	}

	if !checkPassword(req.Password, user.Password) {
		response.BadRequest(c, "账号或密码错误")
		return
	}

	token, err := middleware.GenerateToken(
		user.ID, "student", user.No,
		h.JWTConfig.GetSecret(), h.JWTConfig.GetExpireHour(),
	)
	if err != nil {
		response.InternalError(c, "Token 生成失败")
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "登录成功",
		Data: LoginResponse{
			Token:    token,
			UserID:   user.ID,
			Username: user.No,
			Role:     "student",
			UserInfo: studentInfo(user),
		},
	})
}

// GetCurrentUser GET /auth/me - 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	username, _ := c.Get(middleware.ContextUserName)

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"user_id":  userID,
			"role":     role,
			"username": username,
		},
	})
}

// LecturerRegister POST /auth/lecturer/register - 讲师注册
func (h *AuthHandler) LecturerRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, formatValidationError(err))
		return
	}

	repo := repository.NewLecturerRepo()

	// 检查编号是否已存在
	if existing, _ := repo.FindByNo(req.No); existing != nil {
		response.BadRequest(c, "该编号已被注册")
		return
	}

	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	// 创建讲师记录
	entity := &lecturer.Lecturer{
		No:       req.No,
		Name:     req.Name,
		Password: hashedPwd,
	}
	if err = repo.Create(entity); err != nil {
		response.InternalError(c, "注册失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "注册成功",
		Data:    lecturerInfo(entity),
	})
}

// StudentRegister POST /auth/student/register - 学员注册
func (h *AuthHandler) StudentRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, formatValidationError(err))
		return
	}

	repo := repository.NewStudentRepo()

	// 检查编号是否已存在
	if existing, _ := repo.FindByNo(req.No); existing != nil {
		response.BadRequest(c, "该编号已被注册")
		return
	}

	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	// 创建学员记录
	entity := &student.Student{
		No:       req.No,
		Name:     req.Name,
		Password: hashedPwd,
	}
	if err = repo.Create(entity); err != nil {
		response.InternalError(c, "注册失败: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "注册成功",
		Data:    studentInfo(entity),
	})
}

// ============================================================
//  内部工具函数
// ============================================================

// checkPassword 比对明文密码和 bcrypt 加密后的哈希值
func checkPassword(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

// HashPassword 对明文密码进行 bcrypt 加密 (注册/修改密码时使用)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// lecturerInfo 提取讲师公开信息 (脱敏)
func lecturerInfo(u *lecturer.Lecturer) map[string]interface{} {
	return map[string]interface{}{
		"id":   u.ID,
		"name": u.Name,
		"no":   u.No,
	}
}

// studentInfo 提取学员公开信息 (脱敏)
func studentInfo(u *student.Student) map[string]interface{} {
	return map[string]interface{}{
		"id":   u.ID,
		"name": u.Name,
		"no":   u.No,
	}
}
