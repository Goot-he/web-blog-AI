package repo

import (
	"context"
	"time"

	"core-server/internal/infras/repo/ent"
	"core-server/internal/infras/repo/ent/category"
	"core-server/internal/infras/repo/ent/predicate"
	"core-server/internal/model/entity"

	"entgo.io/ent/dialect/sql"
)

// rootCategoryID 表示一级类型（type_id = 0）
const rootCategoryID uint64 = 0

type CategoryRepo struct {
	*EntClient
}

func NewCategoryRepo(client *EntClient) *CategoryRepo {
	return &CategoryRepo{client}
}

// Create 新建分类（ParentID 为 0 表示一级类型，否则是二级分类）
func (r *CategoryRepo) Create(ctx context.Context, data *entity.Category) error {
	created, err := r.DB(ctx).Category.Create().
		SetCategory(data.Category).
		SetTypeID(data.ParentID).
		Save(ctx)
	if err != nil {
		return err
	}
	data.ID = created.ID
	data.CreatedAt = created.CreatedAt
	data.UpdatedAt = created.UpdatedAt
	return nil
}

// Delete 软删除二级分类
func (r *CategoryRepo) Delete(ctx context.Context, categoryID uint64) error {
	affected, err := r.DB(ctx).Category.Update().
		Where(
			category.IDEQ(categoryID),
			category.TypeIDNEQ(rootCategoryID), // 只允许删二级分类
			category.DeletedAtIsNil(),
		).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrorCategoryNotFound
	}
	return nil
}

// DeleteType 软删除一级类型及其下的所有二级分类
func (r *CategoryRepo) DeleteType(ctx context.Context, typeID uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
		exists, err := r.DB(ctx).Category.Query().
			Where(
				category.IDEQ(typeID),
				category.TypeIDEQ(rootCategoryID),
				category.DeletedAtIsNil(),
			).
			Exist(ctx)
		if err != nil {
			return err
		}
		if !exists {
			return ErrorCategoryNotFound
		}

		now := time.Now()
		// 先删其下的二级分类
		if _, err = r.DB(ctx).Category.Update().
			Where(category.TypeIDEQ(typeID), category.DeletedAtIsNil()).
			SetDeletedAt(now).
			Save(ctx); err != nil {
			return err
		}
		// 再删一级类型本身
		_, err = r.DB(ctx).Category.Update().
			Where(category.IDEQ(typeID), category.DeletedAtIsNil()).
			SetDeletedAt(now).
			Save(ctx)
		return err
	})
}

// ListTypes 查询所有一级类型
func (r *CategoryRepo) ListTypes(ctx context.Context) ([]*entity.Category, error) {
	return r.list(ctx, category.TypeIDEQ(rootCategoryID))
}

// ListByParent 查询某个一级类型下的所有二级分类
func (r *CategoryRepo) ListByParent(ctx context.Context, parentID uint64) ([]*entity.Category, error) {
	return r.list(ctx, category.TypeIDEQ(parentID))
}

// list 查询分类列表（按 ID 升序）
func (r *CategoryRepo) list(ctx context.Context, preds ...predicate.Category) ([]*entity.Category, error) {
	items, err := r.DB(ctx).Category.Query().
		Where(append(preds, category.DeletedAtIsNil())...).
		Order(category.ByID(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityCategories(items), nil
}

// GetByID 查询单个分类
func (r *CategoryRepo) GetByID(ctx context.Context, categoryID uint64) (*entity.Category, error) {
	data, err := r.DB(ctx).Category.Query().
		Where(category.IDEQ(categoryID), category.DeletedAtIsNil()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrorCategoryNotFound
		}
		return nil, err
	}
	return toEntityCategory(data), nil
}
