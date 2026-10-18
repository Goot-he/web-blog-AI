package dto

import "core-server/internal/model/entity"

type CategoryType struct {
	ID   uint64
	Name string
}

type Category struct {
	ID       uint64
	ParentID uint64
	Name     string
}

// ---------- Type ----------

type CreateTypeRequest struct {
	Name string
}
type DeleteTypeRequest struct {
	ID uint64
}
type ListTypesRequest struct{}
type ListTypesResponse struct {
	Types []*CategoryType
}

// ---------- Category ----------

type CreateCategoryRequest struct {
	ParentID uint64
	Name     string
}
type DeleteCategoryRequest struct {
	ID uint64
}
type ListCategoriesRequest struct {
	ParentID uint64
}
type ListCategoriesResponse struct {
	Categories []*Category
}

// ToCategoryTypes entity.Category 列表转一级类型列表
func ToCategoryTypes(items []*entity.Category) []*CategoryType {
	if len(items) == 0 {
		return nil
	}
	out := make([]*CategoryType, 0, len(items))
	for _, item := range items {
		out = append(out, &CategoryType{
			ID:   item.ID,
			Name: item.Category,
		})
	}
	return out
}

// ToCategories entity.Category 列表转二级分类列表
func ToCategories(items []*entity.Category) []*Category {
	if len(items) == 0 {
		return nil
	}
	out := make([]*Category, 0, len(items))
	for _, item := range items {
		out = append(out, &Category{
			ID:       item.ID,
			ParentID: item.ParentID,
			Name:     item.Category,
		})
	}
	return out
}
