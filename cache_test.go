package dingtalk

import (
	"testing"
	"time"
)

type mockExpirable struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	Created   int64  `json:"created"`
}

func (m *mockExpirable) IsExpired() bool {
	return time.Now().Unix() > m.Created+int64(m.ExpiresIn-60)
}

func TestCache(t *testing.T) {
	now := time.Now()
	cases := []struct {
		e      *mockExpirable
		hasErr bool
	}{
		{
			e: &mockExpirable{
				Token:     "token",
				ExpiresIn: 100,
				Created:   now.Unix(),
			},
		},
		{
			e: &mockExpirable{
				Token:     "token",
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

		var n mockExpirable
		err := cache.Get(&n)
		hasErr := err != nil
		if c.hasErr != hasErr {
			t.Fatalf("get cache got %v, expected %v", hasErr, c.hasErr)
		}
		if err != nil && c.e.Token != n.Token {
			t.Fatalf("get cache got %v, expected %v", n.Token, c.e.Token)
		}
	}
}
