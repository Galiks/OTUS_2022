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
		fmt.Printf("Test %q start\n", "empty cache")
		c, err := NewCache(10)
		require.NoError(t, err)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
		fmt.Printf("Test %q done\n", "empty cache")
	})

	t.Run("simple", func(t *testing.T) {
		fmt.Printf("Test %q start\n", "simple")
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
		fmt.Printf("Test %q done\n", "simple")
	})

	t.Run("purge logic (capacity)", func(t *testing.T) {
		fmt.Printf("Test %q start\n", "purge logic (capacity)")
		c, err := NewCache(5)
		require.NoError(t, err)
		require.NotNil(t, c)

		for i := 0; i < 6; i++ {
			wasInCache := c.Set(Key(fmt.Sprint(i)), i+1)
			require.False(t, wasInCache)
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
		fmt.Printf("Test %q done\n", "purge logic (capacity)")
	})

	t.Run("purge logic (timeout)", func(t *testing.T) {
		fmt.Printf("Test %q start\n", "purge logic (timeout)")
		c, err := NewCache(3)
		require.NoError(t, err)
		require.NotNil(t, c)

		for i := 0; i < 3; i++ {
			wasInCache := c.Set(Key(fmt.Sprint(i)), i+1)
			require.False(t, wasInCache)
		}

		for i := 0; i < 5; i++ {
			wasInCache := c.Set(Key(fmt.Sprint(1)), i)
			require.True(t, wasInCache)
			wasInCache = c.Set(Key(fmt.Sprint(2)), i)
			require.True(t, wasInCache)
		}

		wasInCache := c.Set(Key(fmt.Sprint(3)), 4)
		require.False(t, wasInCache)

		val, ok := c.Get(Key(fmt.Sprint(1)))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(fmt.Sprint(2)))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(fmt.Sprint(3)))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(fmt.Sprint(0)))
		require.False(t, ok)
		require.Nil(t, val)
		fmt.Printf("Test %q done\n", "purge logic (timeout)")
	})

	t.Run("capacity error", func(t *testing.T) {
		fmt.Printf("Test %q start\n", "capacity error")
		_, err := NewCache(0)
		require.ErrorIs(t, err, ErrCapacity)
		fmt.Printf("Test %q done\n", "capacity error")
	})
}

func TestCacheMultithreading(t *testing.T) {
	fmt.Printf("Test %q start\n", "Cache Multithreading")
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
	fmt.Printf("Test %q done\n", "Cache Multithreading")
}
