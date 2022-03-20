package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"sync"
	"time"
)

type RSession struct {
	client  *redis.Client
	timeout time.Duration
}

type RMutex struct {
	s   *RSession
	key string
	val string
}

func NewSession(client *redis.Client, timeout time.Duration) *RSession {
	return &RSession{client: client, timeout: timeout}
}

func NewRMutex(s *RSession, key string) *RMutex {
	val := uuid.New().String()
	return &RMutex{s: s, key: key, val: val}
}

func (m *RMutex) Lock(ctx context.Context) error {
	if m.key == "" {
		return errors.New("key empty")
	}
	// NX: 当key不存在时设置值
	args := redis.SetArgs{Mode: "NX", TTL: m.s.timeout}
	cmd := m.s.client.SetArgs(ctx, m.key, m.val, args)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (m *RMutex) UnLock(ctx context.Context) error {
	if m.key == "" {
		return errors.New("key empty")
	}
	cmd := m.s.client.Get(ctx, m.key)
	// 释放锁时对应的值要相等，避免释放错误
	if cmd.Val() != m.val {
		return errors.New(fmt.Sprintf("val not equal: %s %s", cmd.String(), m.val))
	}
	m.s.client.Del(ctx, m.key)
	return nil
}

func initClient() *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       1,
		},
	)
	return client
}

func main() {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       1,
		},
	)
	ctx := context.Background()
	session := NewSession(client, 30*time.Second)

	key := "dist_test"
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		idx := i
		// 并发获取
		go func(idx int) {
			fmt.Println("idx", idx)
			defer wg.Done()
			m := NewRMutex(session, key)
			if err := m.Lock(ctx); err != nil {
				fmt.Println("lock failed: ", idx, err)
				return
			}
			fmt.Println("acquired lock: ", idx)
		}(idx)
	}
	wg.Wait()
}
