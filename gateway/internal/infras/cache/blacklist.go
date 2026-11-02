package cache

import (
	"context"
	"crypto/sha256"
	"fmt"
)

type Blacklist struct {
	*CacheClient
}

func NewBlacklist(client *CacheClient) *Blacklist {
	return &Blacklist{
		CacheClient: client,
	}
}

func (b *Blacklist) IsTokenBlacklist(ctx context.Context, token string) (bool, error) {
	count, err := b.Cache.Exists(ctx, tokenBlacklistKey(token)).Result()
	return count > 0, err
}

func tokenBlacklistKey(token string) string {
	return fmt.Sprintf("blacklist:token:%x", sha256.Sum256([]byte(token)))
}

// ===================================================
// 用户被强制拉黑时的缓存

func (b *Blacklist) AddUser(ctx context.Context, userId uint64) error {
	return b.Cache.Set(ctx, userIdBlacklistKey(userId), "", 0).Err()
}

func (b *Blacklist) RemoveUser(ctx context.Context, userId uint64) error {
	return b.Cache.Del(ctx, userIdBlacklistKey(userId)).Err()
}

func (b *Blacklist) IsUserBlacklisted(ctx context.Context, userid uint64) (bool, error) {
	exits, err := b.Cache.Exists(ctx, userIdBlacklistKey(userid)).Result()
	if err != nil {
		return false, err
	}
	return exits > 0, nil
}

func userIdBlacklistKey(userId uint64) string {
	return fmt.Sprintf("auth:blacklist:ID:%d", userId)
}
