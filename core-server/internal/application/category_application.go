package application

import (
	"context"
	"errors"
	"strings"

	"core-server/internal/domain"
	"core-server/internal/infras/clog"
	"core-server/internal/model/dto"
	"core-server/internal/model/entity"

	"go.uber.org/zap"
)

type CategoryService struct {
	CategoryDomain domain.CategoryRepoDomain
	log            *clog.Log
}

func NewCategoryService(categoryRepo domain.CategoryRepoDomain, log *clog.Log) *CategoryService {
	return &CategoryService{
		CategoryDomain: categoryRepo,
		log:            log,
	}
}

// CreateType 新建一级类型
func (s *CategoryService) CreateType(ctx context.Context, req dto.CreateTypeRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		s.log.Debug("empty category type name")
		return errors.New("type name is required")
	}
	if err := s.CategoryDomain.CreateType(ctx, &entity.Category{
		Category: name,
		ParentID: 0,
	}); err != nil {
		s.log.Error("create category type failed",
			zap.Error(err),
			zap.String("name", name),
		)
		return err
	}
	s.log.Info("category type created", zap.String("name", name))
	return nil
}

// DeleteType 删除一级类型及其下的所有二级分类
func (s *CategoryService) DeleteType(ctx context.Context, req dto.DeleteTypeRequest) error {
	if req.ID == 0 {
		s.log.Debug("empty category type id")
		return errors.New("type id is required")
	}
	if err := s.CategoryDomain.DeleteType(ctx, req.ID); err != nil {
		s.log.Error("delete category type failed",
			zap.Error(err),
			zap.Uint64("type_id", req.ID),
		)
		return err
	}
	s.log.Info("category type deleted", zap.Uint64("type_id", req.ID))
	return nil
}

// ListTypes 查询所有一级类型
func (s *CategoryService) ListTypes(ctx context.Context, _ dto.ListTypesRequest) (*dto.ListTypesResponse, error) {
	items, err := s.CategoryDomain.ListType(ctx)
	if err != nil {
		s.log.Warn("list category types failed", zap.Error(err))
		return nil, err
	}
	s.log.Debug("list category types", zap.Int("count", len(items)))
	return &dto.ListTypesResponse{
		Types: dto.ToCategoryTypes(items),
	}, nil
}

// CreateCategory 在指定一级类型下新建二级分类
func (s *CategoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		s.log.Debug("empty category name", zap.Uint64("parent_id", req.ParentID))
		return errors.New("category name is required")
	}
	if req.ParentID == 0 {
		s.log.Debug("empty parent id", zap.String("name", name))
		return errors.New("parentID is required")
	}
	if err := s.CategoryDomain.CreateCate(ctx, &entity.Category{
		Category: name,
		ParentID: req.ParentID,
	}); err != nil {
		s.log.Error("create category failed",
			zap.Error(err),
			zap.Uint64("parent_id", req.ParentID),
			zap.String("name", name),
		)
		return err
	}
	s.log.Info("category created",
		zap.Uint64("parent_id", req.ParentID),
		zap.String("name", name),
	)
	return nil
}

// DeleteCategory 删除二级分类
func (s *CategoryService) DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) error {
	if req.ID == 0 {
		s.log.Debug("empty category id")
		return errors.New("category id is required")
	}
	if err := s.CategoryDomain.DeleteCate(ctx, req.ID); err != nil {
		s.log.Error("delete category failed",
			zap.Error(err),
			zap.Uint64("category_id", req.ID),
		)
		return err
	}
	s.log.Info("category deleted", zap.Uint64("category_id", req.ID))
	return nil
}

// ListCategories 查询某个一级类型下的二级分类
func (s *CategoryService) ListCategories(ctx context.Context, req dto.ListCategoriesRequest) (*dto.ListCategoriesResponse, error) {
	if req.ParentID == 0 {
		return &dto.ListCategoriesResponse{}, nil
	}
	items, err := s.CategoryDomain.ListCate(ctx, req.ParentID)
	if err != nil {
		s.log.Warn("list categories failed",
			zap.Error(err),
			zap.Uint64("parent_id", req.ParentID),
		)
		return nil, err
	}
	s.log.Debug("list categories",
		zap.Uint64("parent_id", req.ParentID),
		zap.Int("count", len(items)),
	)
	return &dto.ListCategoriesResponse{
		Categories: dto.ToCategories(items),
	}, nil
}
