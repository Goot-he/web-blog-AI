package repo

import (
	"context"
	"core-server/internal/config"
	"core-server/internal/infras/repo/ent"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	// 注册mysql驱动 否则 sql.Open("mysql", dsn) 会报 unknown driver
	_ "github.com/go-sql-driver/mysql"
)

type EntClient struct {
	db     *ent.Client
	driver dialect.Driver
}

// NewEntClient 创建 ent 客户端（连接 MySQL 并设置连接池参数）
func NewEntClient(cfg *config.Config) (*EntClient, error) {
	dsn := cfg.DSN()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)

	return &EntClient{
		db:     ent.NewClient(ent.Driver(drv)),
		driver: drv,
	}, nil
}

// Close 关闭数据库连接
func (c *EntClient) Close() error {
	return c.db.Close()
}

// DB 获取本次操作用的 ent client（在事务中优先返回事务 client，否则用默认连接池）
func (c *EntClient) DB(ctx context.Context) *ent.Client {
	// Transaction contexts carry a transactional Ent client. Normal request
	// contexts do not, so fall back to the base client for regular operations.
	if db := ent.FromContext(ctx); db != nil {
		return db
	}
	return c.db
}

// WithTransaction 在事务中执行 fn（返回错误则回滚，返回 nil 则提交，嵌套调用时复用外层事务）
func (c *EntClient) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	// 已经在事务里了：直接复用，不再开新事务，保证嵌套调用时共用一个事务
	if ent.FromContext(ctx) != nil {
		return fn(ctx)
	}

	tx, err := c.db.Tx(ctx)
	if err != nil {
		return err
	}
	// 兜底：任何没走到Commit的路径（包括fn panic）都会回滚
	defer tx.Rollback()

	// 回滚失败想让它可见的话 就保留这一段； 不想就一行 return err
	if err := fn(ent.NewContext(ctx, tx.Client())); err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf("%w: rolling back transaction: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit()
}
