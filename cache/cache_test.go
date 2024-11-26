package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSetCache(t *testing.T) {
	cache := NewCache()

	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "Value with a type of string",
			key:   "key1",
			value: "value1",
		},
		{
			name:  "Value with a type of integer",
			key:   "key2",
			value: 42,
		},
		{
			name:  "Value with a type of float",
			key:   "key3",
			value: 3.14,
		},
		{
			name:  "Value with a type of boolean",
			key:   "key4",
			value: true,
		},
		{
			name:  "Value with a type of struct",
			key:   "key5",
			value: struct{ Name string }{"Test"},
		},
		{
			name:  "Value with a type of slice",
			key:   "key6",
			value: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache.Set(tc.key, tc.value, time.Minute)

			actualItem, ok := cache.Get(tc.key)
			assert.True(t, ok, "Key should exist in cache")

			assert.Equal(t, tc.value, actualItem, "Actual value should match original value")
		})
	}
}

func TestGetSetCacheJSON(t *testing.T) {
	cache := NewCache()

	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name: "JSON Struct",
			key:  "user1",
			value: struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Active bool   `json:"active"`
			}{
				ID:   1,
				Name: "John Doe",
			},
		},
		{
			name: "JSON Map",
			key:  "config",
			value: map[string]interface{}{
				"database": "mysql",
				"port":     3306,
				"enabled":  true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache.Set(tc.key, tc.value, time.Minute)

			actualItem, ok := cache.Get(tc.key)
			assert.True(t, ok, "Key should exist in cache")

			originalJSON, err := json.Marshal(tc.value)
			assert.NoError(t, err, "Should be able to marshal original value")

			actualJSON, err := json.Marshal(actualItem)
			assert.NoError(t, err, "Should be able to marshal actual value")

			assert.JSONEq(t, string(originalJSON), string(actualJSON))
		})
	}
}

func TestGetExpired(t *testing.T) {
	cache := NewCache()

	tests := []struct {
		name          string
		key           string
		value         interface{}
		ttl           time.Duration
		sleepDuration time.Duration
	}{
		{
			name:          "Get expired key",
			key:           "key1",
			value:         "value1",
			ttl:           time.Millisecond,
			sleepDuration: time.Second,
		},
		{
			name:          "Get not yet expired key",
			key:           "key2",
			value:         "value2",
			ttl:           time.Second * 2,
			sleepDuration: time.Millisecond * 500,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache.Set(tc.key, tc.value, tc.ttl)
			time.Sleep(time.Second)
			actualItem, ok := cache.Get(tc.key)

			if tc.sleepDuration > tc.ttl {
				assert.False(t, ok, "Key should not exist in cache")
				assert.Nil(t, actualItem, "Actual item should be equal to Nil")
			} else {
				assert.True(t, ok, "Key should exist in cache")
				assert.Equal(t, tc.value, actualItem, "Actual value should match original value")
			}
		})
	}
}
