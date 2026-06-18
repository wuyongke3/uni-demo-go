package response

import (
	"net/http"

	"unigo/errorcode"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应结构
type Response struct {
	Code    int         `json:"code"`              // 业务错误码: 0=成功, 其他=错误 (参见 errorcode 包)
	Message string      `json:"message"`           // 提示信息
	Data    interface{} `json:"data,omitempty"`     // 响应数据(可选)
	Details interface{} `json:"details,omitempty"`  // 错误详情(参数校验时返回字段级错误列表)
}

// FieldError 字段级校验错误（用于 details）
type FieldError struct {
	Field   string `json:"field"`   // 出错的字段名
	Message string `json:"message"` // 该字段的错误描述
}

// ============================================================
//  成功响应
// ============================================================

// Success 成功响应 (带数据)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应 (自定义消息)
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success,
		Message: message,
		Data:    data,
	})
}

// ============================================================
//  错误响应 — 核心方法
// ============================================================

// Fail 通用业务错误响应 (HTTP 200 + 业务错误码)
//
// 设计原则：HTTP 状态码表达协议层语义，业务码表达业务层结果。
// 所有业务错误统一返回 HTTP 200，由前端根据 code 判断。
func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errorcode.Msg(code),
	})
}

// FailWithMessage 业务错误响应 (自定义提示信息)
func FailWithMessage(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// FailWithDetails 参数校验错误响应 (附带字段级详情)
func FailWithDetails(c *gin.Context, code int, message string, details []FieldError) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Details: details,
	})
}

// ============================================================
//  HTTP 协议层错误 (真正的 4xx/5xx，用于非业务场景)
// ============================================================

// HTTPError 协议层错误 (如 JSON 解析失败等)
func HTTPError(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}

// ============================================================
//  快捷方法 — 兼容旧调用方式，内部映射到新错误码
// ============================================================

// BadRequest 参数错误 → 映射到 ParamValidateFailed(40001)
func BadRequest(c *gin.Context, message string) {
	FailWithMessage(c, errorcode.ParamValidateFailed, message)
}

// NotFound 资源未找到 → 映射到 ResourceNotFound(30001)
func NotFound(c *gin.Context, message string) {
	FailWithMessage(c, errorcode.ResourceNotFound, message)
}

// InternalError 服务器内部错误 → 映射到 SystemError(10001)
func InternalError(c *gin.Context, message string) {
	FailWithMessage(c, errorcode.SystemError, message)
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页条数
}
