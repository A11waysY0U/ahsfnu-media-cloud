package services

import (
	"sync"
	"time"
)

// UsedNonceStore 记录已使用的 nonce，带过期时间
// 作用：防止同一个 nonce 被重放使用
// 注意：为简化实现，使用进程内存存储；若需要多实例部署，请改为共享存储（如 Redis）
type UsedNonceStore struct {
	mu    sync.Mutex
	store map[string]time.Time // nonce -> expiresAt
	ttl   time.Duration
}

// NewUsedNonceStore 创建一个新的 UsedNonceStore
func NewUsedNonceStore(ttl time.Duration) *UsedNonceStore {
	return &UsedNonceStore{
		store: make(map[string]time.Time),
		ttl:   ttl,
	}
}

// UseOnce 原子性地检查并标记 nonce 已使用
// 返回 true 表示首次使用（合法），返回 false 表示已被使用且仍在有效期内
func (s *UsedNonceStore) UseOnce(nonce string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if expiresAt, exists := s.store[nonce]; exists {
		// 已存在且未过期，判定重放
		if now.Before(expiresAt) {
			return false
		}
		// 过期则允许再次使用，同时刷新有效期
	}

	// 标记为已使用，并设置过期时间
	s.store[nonce] = now.Add(s.ttl)

	// 顺带清理少量过期项，避免无限增长
	if len(s.store) > 1000 {
		s.cleanupLocked()
	}

	return true
}

// cleanupLocked 清理已过期项（调用方需在持有锁时调用）
func (s *UsedNonceStore) cleanupLocked() {
	now := time.Now()
	for k, v := range s.store {
		if now.After(v) {
			delete(s.store, k)
		}
	}
}

// 默认全局实例（15分钟内同一 nonce 只能使用一次）
var defaultUsedNonceStore = NewUsedNonceStore(15 * time.Minute)

// UseNonceOnce 对外暴露的便捷函数
func UseNonceOnce(nonce string) bool {
	return defaultUsedNonceStore.UseOnce(nonce)
}
