package dingtalk

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	now := time.Now()
	cases := []struct {
		e      *KVExpirable
		hasErr bool
	}{
		{
			e: &KVExpirable{
				Key:       "key1",
				Value:     "token",
				ExpiresIn: 100,
				Created:   now.Unix(),
			},
		},
		{
			e: &KVExpirable{
				Key:       "key2",
				Value:     "token",
				ExpiresIn: 100,
				Created:   now.Unix(),
			},
		},
		{
			e: &KVExpirable{
				Key:       "key1",
				Value:     "token",
				ExpiresIn: 10,
				Created:   now.Unix(),
			},
			hasErr: true,
		},
		{
			e: &KVExpirable{
				Key:       "key2",
				Value:     "token",
				ExpiresIn: 10,
				Created:   now.Unix(),
			},
			hasErr: true,
		},
	}

	cache := NewFileCache(".mock_cache_file")
	for _, c := range cases {
		if err := cache.Set(c.e); err != nil {
			t.Fatalf("set cache %v", err)
		}

		n, err := cache.Get(c.e.Key)
		hasErr := err != nil
		if c.hasErr != hasErr {
			t.Fatalf("get cache got %v, expected %v", hasErr, c.hasErr)
		}
		if err != nil && c.e.Value != n.Value {
			t.Fatalf("get cache got %v, expected %v", n.Value, c.e.Value)
		}
	}

	_, err := cache.Get("nokey")
	if err == nil {
		t.Fatal("should be get an error")
	}
}
