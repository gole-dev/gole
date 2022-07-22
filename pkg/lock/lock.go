package lock

import (
	"context"
	"time"
)

const (
	// RedisLockKey redis lock key
	RedisLockKey = "gole:redis:lock:%s"
	// EtcdLockKey etcd lock key
	EtcdLockKey = "/gole/lock/%s"
)

// Lock define common func
type Lock interface {
	Lock(ctx context.Context, timeout time.Duration) (bool, error)
	Unlock(ctx context.Context) (bool, error)
}
