package gocash

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache(CacheOptions{DefaultTimeout: 1 * time.Second})
	if c == nil {
		t.Fatal("NewCache returned nil")
	}
}

func TestGetMissingKey(t *testing.T) {
	c := NewCache(CacheOptions{DefaultTimeout: 10 * time.Second})
}
