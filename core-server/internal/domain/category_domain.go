package domain

import (
	"context"

	"core-server/internal/model/entity"
)

type CategoryRepo interface {
	// Create 新建分类（ParentID 为 0 表示一级类型，否则是二级分类）
	Create(ctx context.Context, category *entity.Category) error
	// Delete 软删除二级分类
	Delete(ctx context.Context, categoryID uint64) error
	// DeleteType 软删除一级类型及其下的所有二级分类
	DeleteType(ctx context.Context, typeID uint64) error
	// ListTypes 查询所有一级类型
	ListTypes(ctx context.Context) ([]*entity.Category, error)
	// ListByParent 查询某个一级类型下的所有二级分类
	ListByParent(ctx context.Context, parentID uint64) ([]*entity.Category, error)
	// GetByID 查询单个分类
	GetByID(ctx context.Context, categoryID uint64) (*entity.Category, error)
}

type CategoryDomainRepo interface {
	CategoryRepo
	ITransaction
}
