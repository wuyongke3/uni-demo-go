package service

import (
	"errors"
	"unigo/repository"

	"gorm.io/gorm"
)

// Operations 通用 CRUD 操作接口 (Handler 层依赖此接口)
type Operations[T any] interface {
	Create(entity *T) (*T, error)
	GetByID(id uint) (*T, error)
	GetByIDs(ids []uint) ([]T, error)
	List(page, pageSize int) ([]T, int64, error)
	Update(id uint, entity *T) (*T, error)
	Delete(id uint) error
	BatchDelete(ids []uint) error
}

// CRUDService 泛型 CRUD 服务基类 (封装通用业务逻辑)
type CRUDService[T any, R repository.GenericRepository[T]] struct {
	repo R
}

// NewCRUDService 创建泛型服务实例
func NewCRUDService[T any, R repository.GenericRepository[T]](repo R) *CRUDService[T, R] {
	return &CRUDService[T, R]{repo: repo}
}

// Create 创建记录
func (s *CRUDService[T, R]) Create(entity *T) (*T, error) {
	if err := s.repo.Create(entity); err != nil {
		return nil, errors.New("创建失败: " + err.Error())
	}
	return entity, nil
}

// GetByID 根据ID查询
func (s *CRUDService[T, R]) GetByID(id uint) (*T, error) {
	e, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("记录不存在")
		}
		return nil, err
	}
	return e, nil
}

// GetByIDs 批量根据ID列表查询
func (s *CRUDService[T, R]) GetByIDs(ids []uint) ([]T, error) {
	return s.repo.GetByIDs(ids)
}

// List 分页查询
func (s *CRUDService[T, R]) List(page, pageSize int) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

// Update 更新记录
func (s *CRUDService[T, R]) Update(id uint, entity *T) (*T, error) {
	if _, err := s.repo.GetByID(id); err != nil {
		return nil, err
	}
	if err := s.repo.Update(id, entity); err != nil {
		return nil, errors.New("更新失败: " + err.Error())
	}
	return entity, nil
}

// Delete 删除单条记录
func (s *CRUDService[T, R]) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// BatchDelete 批量删除
func (s *CRUDService[T, R]) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("ID 列表不能为空")
	}
	return s.repo.BatchDelete(ids)
}
