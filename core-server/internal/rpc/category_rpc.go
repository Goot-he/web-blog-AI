package rpc

import (
	"context"

	"core-server/internal/application"
	"core-server/internal/model/dto"
	"core-server/internal/rpc/categorypb"
)

type CategoryRpc struct {
	categorypb.UnimplementedCategoryServiceServer
	categoryService *application.CategoryService
}

func NewCategoryRpc(categoryService *application.CategoryService) *CategoryRpc {
	return &CategoryRpc{
		categoryService: categoryService,
	}
}

// CreateType 新建一级类型
func (r *CategoryRpc) CreateType(ctx context.Context, req *categorypb.CreateTypeRequest) (*categorypb.CreateTypeResponse, error) {
	err := r.categoryService.CreateType(ctx, dto.CreateTypeRequest{
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return &categorypb.CreateTypeResponse{Success: true}, nil
}

// DeleteType 删除一级类型及其子分类
func (r *CategoryRpc) DeleteType(ctx context.Context, req *categorypb.DeleteTypeRequest) (*categorypb.DeleteTypeResponse, error) {
	err := r.categoryService.DeleteType(ctx, dto.DeleteTypeRequest{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &categorypb.DeleteTypeResponse{Success: true}, nil
}

// ListTypes 查询所有一级类型
func (r *CategoryRpc) ListTypes(ctx context.Context, _ *categorypb.ListTypesRequest) (*categorypb.ListTypesResponse, error) {
	resp, err := r.categoryService.ListTypes(ctx, dto.ListTypesRequest{})
	if err != nil {
		return nil, err
	}
	types := make([]*categorypb.CategoryType, 0, len(resp.Types))
	for _, t := range resp.Types {
		types = append(types, &categorypb.CategoryType{
			Id:   t.ID,
			Name: t.Name,
		})
	}
	return &categorypb.ListTypesResponse{Types: types}, nil
}

// CreateCategory 在指定一级类型下新建二级分类
func (r *CategoryRpc) CreateCategory(ctx context.Context, req *categorypb.CreateCategoryRequest) (*categorypb.CreateCategoryResponse, error) {
	err := r.categoryService.CreateCategory(ctx, dto.CreateCategoryRequest{
		ParentID: req.GetParentID(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, err
	}
	return &categorypb.CreateCategoryResponse{Success: true}, nil
}

// DeleteCategory 删除二级分类
func (r *CategoryRpc) DeleteCategory(ctx context.Context, req *categorypb.DeleteCategoryRequest) (*categorypb.DeleteCategoryResponse, error) {
	err := r.categoryService.DeleteCategory(ctx, dto.DeleteCategoryRequest{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &categorypb.DeleteCategoryResponse{Success: true}, nil
}

// ListCategories 查询某个一级类型下的二级分类
func (r *CategoryRpc) ListCategories(ctx context.Context, req *categorypb.ListCategoriesRequest) (*categorypb.ListCategoriesResponse, error) {
	resp, err := r.categoryService.ListCategories(ctx, dto.ListCategoriesRequest{
		ParentID: req.GetParentID(),
	})
	if err != nil {
		return nil, err
	}
	categories := make([]*categorypb.Category, 0, len(resp.Categories))
	for _, c := range resp.Categories {
		categories = append(categories, &categorypb.Category{
			Id:       c.ID,
			ParentID: c.ParentID,
			Name:     c.Name,
		})
	}
	return &categorypb.ListCategoriesResponse{Categories: categories}, nil
}
