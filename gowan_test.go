package gowan

import (
	"testing"
	"time"
)

func TestExpires(t *testing.T) {
	now := time.Now().UTC()

	cases := []struct {
		expected bool
		item     *item
		time     int64
	}{
		{
			expected: true,
			item: &item{
				expires: now.Add(-5 * time.Minute).UnixNano(),
			},
			time: now.UnixNano(),
		},
		{
			expected: true,
			item: &item{
				expires: 0,
			},
			time: now.UnixNano(),
		},
		{
			expected: false,
			item: &item{
				expires: now.Add(5 * time.Minute).UnixNano(),
			},
			time: now.UnixNano(),
		},
	}

	for _, c := range cases {
		var actual = c.item.Expired(c.time)
		if actual != c.expected {
			t.Errorf("result = %t expected. got result = %t", c.expected, actual)
		}
	}
}

func TestPutAndGet(t *testing.T) {
	now := time.Now().UTC()

	cases := []struct {
		expected interface{}
		key      string
		value    interface{}
		expires  int64
	}{
		{
			expected: "cached-value",
			key:      "cached-key",
			value:    "cached-value",
			expires:  now.Add(10 * time.Minute).UnixNano(),
		},
		{
			// test case for overwriting a item.
			expected: "overwrite-cached-value",
			key:      "cached-key",
			value:    "overwrite-cached-value",
			expires:  now.Add(10 * time.Minute).UnixNano(),
		},
		{
			// test case for an expired item.
			expected: nil,
			key:      "expired-key",
			value:    "expired-value",
			expires:  now.Add(-10 * time.Minute).UnixNano(),
		},
		{
			expected: 1,
			key:      "int-cached-key",
			value:    1,
			expires:  now.Add(10 * time.Minute).UnixNano(),
		},
		{
			expected: struct{ value string }{value: "struct"},
			key:      "struct-cached-key",
			value:    struct{ value string }{value: "struct"},
			expires:  now.Add(10 * time.Minute).UnixNano(),
		},
	}

	cache := New()
	for _, c := range cases {
		cache.Put(c.key, c.value, c.expires)
		time.Sleep(time.Second * 1)
		var actual = cache.Get(c.key)
		if actual != c.expected {
			t.Errorf("result = %v expected. got result = %v", c.expected, actual)
		}
	}
}
