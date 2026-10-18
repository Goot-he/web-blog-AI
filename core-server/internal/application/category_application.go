package application

import (
	"context"
	"errors"
	"strings"

	"core-server/internal/domain"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"
)

type CategoryService struct {
	CategoryDomain domain.CategoryDomainRepo
}

func NewCategoryService(categoryRepo domain.CategoryDomainRepo) *CategoryService {
	return &CategoryService{
		CategoryDomain: categoryRepo,
	}
}

// CreateType 新建一级类型
func (s *CategoryService) CreateType(ctx context.Context, req dto.CreateTypeRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("type name is required")
	}
	return s.CategoryDomain.Create(ctx, &entity.Category{
		Category: name,
		ParentID: 0,
	})
}

// DeleteType 删除一级类型及其下的所有二级分类
func (s *CategoryService) DeleteType(ctx context.Context, req dto.DeleteTypeRequest) error {
	if req.ID == 0 {
		return errors.New("type id is required")
	}
	return s.CategoryDomain.DeleteType(ctx, req.ID)
}

// ListTypes 查询所有一级类型
func (s *CategoryService) ListTypes(ctx context.Context, _ dto.ListTypesRequest) (*dto.ListTypesResponse, error) {
	items, err := s.CategoryDomain.ListTypes(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.ListTypesResponse{
		Types: dto.ToCategoryTypes(items),
	}, nil
}

// CreateCategory 在指定一级类型下新建二级分类
func (s *CategoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("category name is required")
	}
	if req.ParentID == 0 {
		return errors.New("parentID is required")
	}

	parent, err := s.CategoryDomain.GetByID(ctx, req.ParentID)
	if err != nil {
		return err
	}
	if parent.ParentID != 0 {
		return errors.New("parent category must be a first-level type")
	}

	return s.CategoryDomain.Create(ctx, &entity.Category{
		Category: name,
		ParentID: req.ParentID,
	})
}

// DeleteCategory 删除二级分类
func (s *CategoryService) DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) error {
	if req.ID == 0 {
		return errors.New("category id is required")
	}
	return s.CategoryDomain.Delete(ctx, req.ID)
}

// ListCategories 查询某个一级类型下的二级分类
func (s *CategoryService) ListCategories(ctx context.Context, req dto.ListCategoriesRequest) (*dto.ListCategoriesResponse, error) {
	if req.ParentID == 0 {
		return &dto.ListCategoriesResponse{}, nil
	}
	items, err := s.CategoryDomain.ListByParent(ctx, req.ParentID)
	if err != nil {
		return nil, err
	}
	return &dto.ListCategoriesResponse{
		Categories: dto.ToCategories(items),
	}, nil
}
