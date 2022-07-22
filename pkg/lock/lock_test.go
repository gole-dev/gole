package lock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gole-dev/gole/pkg/redis"
)

func TestLockWithDefaultTimeout(t *testing.T) {
	redis.InitTestRedis()

	lock := NewRedisLock(redis.RedisClient, "lock1")
	ok, err := lock.Lock(context.Background(), 2*time.Second)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatal("lock is not ok")
	}

	ok, err = lock.Unlock(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Unlock is not ok")
	}

	t.Log(ok)
}

func TestLockWithTimeout(t *testing.T) {
	redis.InitTestRedis()

	t.Run("should lock/unlock success", func(t *testing.T) {
		ctx := context.Background()
		lock1 := NewRedisLock(redis.RedisClient, "lock2")
		ok, err := lock1.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = lock1.Unlock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("should unlock failed", func(t *testing.T) {
		ctx := context.Background()
		lock2 := NewRedisLock(redis.RedisClient, "lock3")
		ok, err := lock2.Lock(ctx, 2*time.Second)
		assert.Nil(t, err)
		assert.True(t, ok)

		time.Sleep(3 * time.Second)

		ok, err = lock2.Unlock(ctx)
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}
