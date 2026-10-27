package repo

import (
	"context"
	"time"

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

// CreateType 新建一级类型
func (r *CategoryRepo) CreateType(ctx context.Context, cate *entity.Category) error {
	return r.create(ctx, cate, rootCategoryID)
}

// CreateCate 新建二级分类（父级必须是一级类型）
func (r *CategoryRepo) CreateCate(ctx context.Context, cate *entity.Category) error {
	if cate.ParentID == rootCategoryID {
		return ErrorInvalidParam
	}
	exists, err := r.DB(ctx).Category.Query().
		Where(
			category.IDEQ(cate.ParentID),
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
	return r.create(ctx, cate, cate.ParentID)
}

// create 插入分类并把生成的字段回填给调用方
func (r *CategoryRepo) create(ctx context.Context, cate *entity.Category, parentID uint64) error {
	created, err := r.DB(ctx).Category.Create().
		SetCategory(cate.Category).
		SetTypeID(parentID).
		Save(ctx)
	if err != nil {
		return err
	}
	cate.ID = created.ID
	cate.ParentID = parentID
	cate.CreatedAt = created.CreatedAt
	cate.UpdatedAt = created.UpdatedAt
	return nil
}

// DeleteType 软删除一级类型及其下的所有二级分类
func (r *CategoryRepo) DeleteType(ctx context.Context, id uint64) error {
	return r.WithTransaction(ctx, func(ctx context.Context) error {
		exists, err := r.DB(ctx).Category.Query().
			Where(
				category.IDEQ(id),
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
			Where(category.TypeIDEQ(id), category.DeletedAtIsNil()).
			SetDeletedAt(now).
			Save(ctx); err != nil {
			return err
		}
		// 再删一级类型本身
		_, err = r.DB(ctx).Category.Update().
			Where(category.IDEQ(id), category.DeletedAtIsNil()).
			SetDeletedAt(now).
			Save(ctx)
		return err
	})
}

// DeleteCate 软删除二级分类
func (r *CategoryRepo) DeleteCate(ctx context.Context, id uint64) error {
	affected, err := r.DB(ctx).Category.Update().
		Where(
			category.IDEQ(id),
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

// ListType 查询所有一级类型
func (r *CategoryRepo) ListType(ctx context.Context) ([]*entity.Category, error) {
	return r.list(ctx, category.TypeIDEQ(rootCategoryID))
}

// ListCate 查询某个一级类型下的所有二级分类
func (r *CategoryRepo) ListCate(ctx context.Context, parentID uint64) ([]*entity.Category, error) {
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
