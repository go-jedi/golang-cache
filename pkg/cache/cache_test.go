package cache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(10 * time.Second)
	defer cache.StartCleanup(1 * time.Second)
}

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("key", "value")

	var value string
	if !cache.Get("key", &value) {
		t.Errorf("failed to get cache item")
	}

	if value != "value" {
		t.Errorf("expected 'value', got '%s'", value)
	}
}

func TestCache_SetWithCustomTTL(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("key", "value", 5*time.Second)

	var value string
	if !cache.Get("key", &value) {
		t.Errorf("failed to get cache item")
	}

	if value != "value" {
		t.Errorf("expected 'value', got '%s'", value)
	}

	time.Sleep(6 * time.Second)

	if !cache.Expired("key") {
		t.Error("cache item did not expire as expected")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(10 * time.Second)

	cache.Set("key", "value")

	cache.Delete("key")

	var value string
	if cache.Get("key", &value) {
		t.Error("cache item was not deleted")
	}
}

func TestCache_Expired(t *testing.T) {
	cache := NewCache(1 * time.Second)

	cache.Set("key", "value")

	time.Sleep(2 * time.Second)

	if !cache.Expired("key") {
		t.Errorf("cache item did not expire as expected")
	}
}

func TestCache_Cleanup(t *testing.T) {
	cache := NewCache(1 * time.Second)
	cache.StartCleanup(1 * time.Second)

	cache.Set("key", "value")

	time.Sleep(2 * time.Second)

	if !cache.Expired("key") {
		t.Errorf("cache item was not cleaned up")
	}
}

type testStruct struct {
	Field string
}

func TestCache_JSONMarshaling(t *testing.T) {
	cache := NewCache(10 * time.Second)

	testData := testStruct{Field: "value"}
	cache.Set("key", testData)

	var retrieved testStruct
	if !cache.Get("key", &retrieved) {
		t.Errorf("failed to get cache item")
	}
	if retrieved.Field != "value" {
		t.Errorf("expected 'value', got '%s'", retrieved.Field)
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(10 * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key string, value interface{}) {
			defer wg.Done()
			cache.Set(key, value)
			var retrieved interface{}
			if !cache.Get(key, &retrieved) {
				t.Errorf("failed to get cache item")
			}
			if retrieved != value {
				t.Errorf("expected '%v', got '%v'", value, retrieved)
			}
		}(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
	}
	wg.Wait()
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache(time.Minute)
	key := "benchmark-key"
	data := "benchmark-data"

	for n := 0; n < b.N; n++ {
		cache.Set(key, data)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache(time.Minute)
	key := "benchmark-key"
	data := "benchmark-data"
	cache.Set(key, data)

	for n := 0; n < b.N; n++ {
		var result string
		cache.Get(key, &result)
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cache := NewCache(time.Minute)
	key := "benchmark-key"
	data := "benchmark-data"
	cache.Set(key, data)

	for n := 0; n < b.N; n++ {
		cache.Delete(key)
		cache.Set(key, data)
	}
}

func BenchmarkCacheExpired(b *testing.B) {
	cache := NewCache(time.Minute)
	key := "benchmark-key"
	data := "benchmark-data"
	cache.Set(key, data)

	for n := 0; n < b.N; n++ {
		cache.Expired(key)
	}
}

func BenchmarkCacheCleanup(b *testing.B) {
	cache := NewCache(time.Second) // Short TTL for frequent cleanup
	key := "benchmark-key"
	data := "benchmark-data"

	for i := 0; i < 100; i++ {
		cache.Set(key+strconv.Itoa(i), data, time.Millisecond)
	}

	for n := 0; n < b.N; n++ {
		cache.Cleanup()
	}
}

func BenchmarkCacheSetMultipleInputs(b *testing.B) {
	cache := NewCache(time.Minute)
	inputs := []string{"key1", "key2", "key3", "key4", "key5"}
	data := "benchmark-data"

	for _, key := range inputs {
		b.Run(key, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				cache.Set(key, data)
			}
		})
	}
}
