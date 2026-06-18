package handler

import (
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unigo/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CRUDOperations 通用 CRUD 操作接口
type CRUDOperations[T any] interface {
	Create(entity *T) (*T, error)
	GetByID(id uint) (*T, error)
	GetByIDs(ids []uint) ([]T, error)
	List(page, pageSize int) ([]T, int64, error)
	Update(id uint, entity *T) (*T, error)
	Delete(id uint) error
	BatchDelete(ids []uint) error
}

// CRUDEntity 通用处理器
type CRUDEntity[T any] struct {
	SVC CRUDOperations[T]
}

// AllRequest 分页查询请求参数
type AllRequest struct {
	Limit int    `json:"limit"` // 每页条数 (默认 10)
	Page  int    `json:"page"`  // 页码 (默认 1, start 优先时忽略)
	Start int    `json:"start"` // 起始偏移量 (优先级高于 page)
	Sort  string `json:"sort"`  // 排序字段 (如 "-name" 表示倒序)
}

// All POST /all - 分页列表查询 (支持排序)
func (h *CRUDEntity[T]) All(c *gin.Context) {
	var req AllRequest
	// 支持 JSON Body 或 Query 参数
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			req = parseQueryParams(c)
		}
	} else {
		req = parseQueryParams(c)
	}

	// 默认值处理
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Start < 0 {
		req.Start = 0
	}

	// start 优先级高于 page 计算页码
	page := req.Page
	if req.Start > 0 && req.Page == 1 {
		page = req.Start/req.Limit + 1
	}
	if page < 1 {
		page = 1
	}

	list, total, err := h.SVC.List(page, req.Limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 排序处理
	if req.Sort != "" {
		list = applySort(list, req.Sort)
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Info GET /info/:ids - 查询详情 (支持多 ID, 带关联数据)
func (h *CRUDEntity[T]) Info(c *gin.Context) {
	idsStr := c.Param("ids")
	ids := parseIDString(idsStr)

	if len(ids) == 0 {
		response.BadRequest(c, "无效的 ID 参数")
		return
	}

	var data interface{}
	var err error

	if len(ids) == 1 {
		// 单条查询
		data, err = h.SVC.GetByID(ids[0])
	} else {
		// 批量查询
		data, err = h.SVC.GetByIDs(ids)
	}

	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, data)
}

// Add POST /add - 新增记录
func (h *CRUDEntity[T]) Add(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}
	result, err := h.SVC.Create(&entity)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "新增成功", result)
}

// Modify POST /modify/:id - 编辑记录
func (h *CRUDEntity[T]) Modify(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		code, msg, details := formatValidationError(err)
		if details != nil {
			response.FailWithDetails(c, code, msg, details)
		} else {
			response.FailWithMessage(c, code, msg)
		}
		return
	}
	result, err := h.SVC.Update(id, &entity)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessWithMessage(c, "编辑成功", result)
}

// Delete DELETE /delete/:ids - 批量删除
func (h *CRUDEntity[T]) Delete(c *gin.Context) {
	idsStr := c.Param("ids")
	ids := parseIDString(idsStr)

	if len(ids) == 0 {
		response.BadRequest(c, "无效的 ID 参数")
		return
	}

	if err := h.SVC.BatchDelete(ids); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	msg := "删除成功"
	if len(ids) > 1 {
		msg = "批量删除成功 (" + strconv.Itoa(len(ids)) + " 条)"
	}
	c.JSON(http.StatusOK, response.Response{Code: 0, Message: msg})
}

// ============================================================
//  内部工具函数
// ============================================================

// parseUintParam 解析路径参数为 uint
func parseUintParam(c *gin.Context, key string) (uint, error) {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// parseIDString 解析逗号分隔的 ID 字符串为 uint 列表
func parseIDString(s string) []uint {
	var ids []uint
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	return ids
}

// parseQueryParams 从 Query 参数解析 AllRequest
func parseQueryParams(c *gin.Context) AllRequest {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	sort := c.Query("sort")
	return AllRequest{Limit: limit, Page: page, Start: start, Sort: sort}
}

// applySort 对切片按字段排序 (支持 "-" 前缀倒序)
func applySort[T any](list []T, sortField string) []T {
	if len(list) <= 1 || sortField == "" {
		return list
	}

	desc := strings.HasPrefix(sortField, "-")
	field := strings.TrimPrefix(sortField, "-")

	result := make([]T, len(list))
	copy(result, list)

	sort.Slice(result, func(i, j int) bool {
		vi := reflect.ValueOf(result[i])
		vj := reflect.ValueOf(result[j])

		// 处理指针类型
		if vi.Kind() == reflect.Ptr {
			vi = vi.Elem()
		}
		if vj.Kind() == reflect.Ptr {
			vj = vj.Elem()
		}

		fi := vi.FieldByName(field)
		fj := vj.FieldByName(field)

		if !fi.IsValid() || !fj.IsValid() {
			return false
		}

		cmp := compareValues(fi, fj)
		if desc {
			return cmp > 0
		}
		return cmp < 0
	})

	return result
}

// compareValues 比较两个 reflect.Value
func compareValues(a, b reflect.Value) int {
	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai, bi := a.Int(), b.Int()
		if ai < bi {
			return -1
		} else if ai > bi {
			return 1
		}
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ai, bi := a.Uint(), b.Uint()
		if ai < bi {
			return -1
		} else if ai > bi {
			return 1
		}
		return 0
	case reflect.Float32, reflect.Float64:
		ai, bi := a.Float(), b.Float()
		if ai < bi {
			return -1
		} else if ai > bi {
			return 1
		}
		return 0
	case reflect.String:
		ai, bi := a.String(), b.String()
		if ai < bi {
			return -1
		} else if ai > bi {
			return 1
		}
		return 0
	default:
		return 0
	}
}

// formatValidationError 格式化校验错误为结构化字段级错误
//
// 返回值: (业务错误码, 概要消息, 字段错误详情列表)
func formatValidationError(err error) (int, string, []response.FieldError) {
	if ve, ok := err.(validator.ValidationErrors); ok {
		details := make([]response.FieldError, 0, len(ve))
		for _, e := range ve {
			fe := response.FieldError{Field: e.Field()}
			switch e.Tag() {
			case "required":
				fe.Message = "该字段为必填项"
			case "min":
				if _, isNum := e.Value().(float64); isNum || e.Type().Kind() == reflect.Int || e.Type().Kind() == reflect.Int64 {
					fe.Message = "最小值为 " + e.Param()
				} else {
					fe.Message = "最小长度为 " + e.Param() + " 个字符"
				}
			case "max":
				if _, isNum := e.Value().(float64); isNum || e.Type().Kind() == reflect.Int || e.Type().Kind() == reflect.Int64 {
					fe.Message = "最大值为 " + e.Param()
				} else {
					fe.Message = "最大长度为 " + e.Param() + " 个字符"
				}
			case "oneof":
				fe.Message = "必须是以下值之一: " + e.Param()
			case "url":
				fe.Message = "必须填写合法的 URL"
			case "email":
				fe.Message = "邮箱格式不正确"
			default:
				fe.Message = "校验失败: " + e.Tag()
			}
			details = append(details, fe)
		}
		return 40001, "参数校验失败", details
	}

	// 非 validator 错误 (如 JSON 解析失败)
	return 40003, "请求参数格式错误: " + err.Error(), nil
}
