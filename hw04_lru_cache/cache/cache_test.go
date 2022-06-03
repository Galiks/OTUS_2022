package cache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c, err := NewCache(10)
		require.NoError(t, err)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5)
		require.NoError(t, err)
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c, err := NewCache(5)
		require.NoError(t, err)
		require.NotNil(t, c)

		for i := 0; i < 6; i++ {
			c.Set(Key(fmt.Sprint(i)), i+1)
		}

		for i := 0; i < 6; i++ {
			if i == 0 {
				val, ok := c.Get(Key(fmt.Sprint(i)))
				require.False(t, ok)
				require.Nil(t, val)
			} else {
				val, ok := c.Get(Key(fmt.Sprint(i)))
				require.True(t, ok)
				require.NotNil(t, val)
			}
		}
	})

	t.Run("capacity error", func(t *testing.T) {
		_, err := NewCache(0)
		require.ErrorIs(t, err, ErrCapacity)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c, _ := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
