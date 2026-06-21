package handler

import (
	"fmt"

	"unigo/database"
	"unigo/errorcode"
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
	Name     string `json:"name" binding:"required,max=50"`    // 姓名 (登录账号)
	Password string `json:"password" binding:"required,min=6"` // 密码
}

// RegisterRequest 注册请求体 (通用)
type RegisterRequest struct {
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
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}

	repo := repository.NewLecturerRepo()
	user, err := repo.FindByName(req.Name)
	if err != nil || user == nil {
		response.Fail(c, errorcode.LoginFailed)
		return
	}

	if !checkPassword(req.Password, user.Password) {
		response.Fail(c, errorcode.LoginFailed)
		return
	}

	token, err := middleware.GenerateToken(
		user.ID, "lecturer", user.Name,
		h.JWTConfig.GetSecret(), h.JWTConfig.GetExpireHour(),
	)
	if err != nil {
		response.FailWithMessage(c, errorcode.TokenGenFailed, "Token 生成失败")
		return
	}

	response.SuccessWithMessage(c, "登录成功", LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Name,
		Role:     "lecturer",
		UserInfo: lecturerInfo(user),
	})
}

// StudentLogin POST /auth/student/login - 学员登录
func (h *AuthHandler) StudentLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}

	repo := repository.NewStudentRepo()
	user, err := repo.FindByName(req.Name)
	if err != nil || user == nil {
		response.Fail(c, errorcode.LoginFailed)
		return
	}

	if !checkPassword(req.Password, user.Password) {
		response.Fail(c, errorcode.LoginFailed)
		return
	}

	token, err := middleware.GenerateToken(
		user.ID, "student", user.Name,
		h.JWTConfig.GetSecret(), h.JWTConfig.GetExpireHour(),
	)
	if err != nil {
		response.FailWithMessage(c, errorcode.TokenGenFailed, "Token 生成失败")
		return
	}

	response.SuccessWithMessage(c, "登录成功", LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Name,
		Role:     "student",
		UserInfo: studentInfo(user),
	})
}

// GetCurrentUser GET /auth/me - 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	username, _ := c.Get(middleware.ContextUserName)

	response.Success(c, gin.H{
		"user_id":  userID,
		"role":     role,
		"username": username,
	})
}

// LecturerRegister POST /auth/lecturer/register - 讲师注册
func (h *AuthHandler) LecturerRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}

	repo := repository.NewLecturerRepo()

	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		response.FailWithMessage(c, errorcode.EncryptFailed, "密码处理失败")
		return
	}

	// 先创建记录获取自增 ID，再生成编号
	entity := &lecturer.Lecturer{
		Name:     req.Name,
		Password: hashedPwd,
	}
	if err = repo.Create(entity); err != nil {
		response.FailWithMessage(c, errorcode.RegisterFailed, "注册失败: "+err.Error())
		return
	}

	// 生成编号: T + 8位序号 (如 T00000001)
	entity.No = generateNo("T", entity.ID)
	database.DB.Model(entity).Update("no", entity.No)

	response.SuccessWithMessage(c, "注册成功", lecturerInfo(entity))
}

// StudentRegister POST /auth/student/register - 学员注册
func (h *AuthHandler) StudentRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}

	repo := repository.NewStudentRepo()

	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		response.FailWithMessage(c, errorcode.EncryptFailed, "密码处理失败")
		return
	}

	// 先创建记录获取自增 ID，再生成编号
	entity := &student.Student{
		Name:     req.Name,
		Password: hashedPwd,
	}
	if err = repo.Create(entity); err != nil {
		response.FailWithMessage(c, errorcode.RegisterFailed, "注册失败: "+err.Error())
		return
	}

	// 生成编号: S + 8位序号 (如 S00000001)
	entity.No = generateNo("S", entity.ID)
	database.DB.Model(entity).Update("no", entity.No)

	response.SuccessWithMessage(c, "注册成功", studentInfo(entity))
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

// generateNo 生成编号: 前缀 + 8位序号 (如 T00000001, S00000042)
func generateNo(prefix string, id uint) string {
	return fmt.Sprintf("%s%08d", prefix, id)
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
