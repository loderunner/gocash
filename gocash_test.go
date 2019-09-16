package gocash

import (
	"strconv"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewCache(CacheOptions{})
	if c == nil {
		t.Fatal("NewCache returned nil")
	}
}

func TestGetMissingKey(t *testing.T) {
	c := NewCache(CacheOptions{})
	w, _ := c.Get("hello")
	if w != nil {
		t.Fatalf("unexpected value for missing key: expected <nil>, got %#v", w)
	}
}

func TestSetGet(t *testing.T) {
	c := NewCache(CacheOptions{})
	d := c.Set("hello", "world")
	if d != NeverExpires {
		t.Errorf("unexpected deadline: expected NeverExpires, received  %s", d)
	}
	v, d2 := c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}
}

func TestSetHas(t *testing.T) {
	c := NewCache(CacheOptions{})
	d := c.Set("hello", "world")
	if d != NeverExpires {
		t.Errorf("unexpected deadline: expected NeverExpires, received  %s", d)
	}
	has, d2 := c.Has("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if !has {
		t.Fatalf("unexpected value: expected has to be true")
	}
}

func TestSetWithTimeout(t *testing.T) {
	c := NewCache(CacheOptions{})
	d := c.SetWithTimeout("hello", "world", 100*time.Millisecond)
	expectedDeadline := time.Now().Add(100 * time.Millisecond)
	if d.After(expectedDeadline) {
		t.Errorf(
			"unexpected deadline: expected > %s, received  %s",
			expectedDeadline,
			d,
		)
	}
	v, d2 := c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}
}

func TestSetWithDeadline(t *testing.T) {
	c := NewCache(CacheOptions{})
	d := time.Now().Add(100 * time.Millisecond)
	d2 := c.SetWithDeadline("hello", "world", d)
	if d != d2 {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	v, d2 := c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}
}

func TestDefaultTimeout(t *testing.T) {
	// Set
	c := NewCache(CacheOptions{DefaultTimeout: 100 * time.Millisecond})
	d := c.Set("hello", "world")
	expectedDeadline := time.Now().Add(100 * time.Millisecond)
	if d.After(expectedDeadline) {
		t.Errorf(
			"unexpected deadline: expected > %s, received  %s",
			expectedDeadline,
			d,
		)
	}
	v, d2 := c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}

	// SetWithTimeout
	d = c.SetWithTimeout("hello", "world", 200*time.Millisecond)
	expectedDeadline = time.Now().Add(200 * time.Millisecond)
	if d.After(expectedDeadline) {
		t.Errorf(
			"unexpected deadline: expected > %s, received  %s",
			expectedDeadline,
			d,
		)
	}
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}

	// SetWithDeadline
	expectedDeadline = time.Now().Add(200 * time.Millisecond)
	d = c.SetWithDeadline("hello", "world", expectedDeadline)
	if d != expectedDeadline {
		t.Errorf(
			"unexpected deadline: expected > %s, received  %s",
			expectedDeadline,
			d,
		)
	}
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != "world" {
		t.Fatalf("unexpected value: expected %#v, received %#v", "world", v)
	}

	time.Sleep(100 * time.Millisecond)
	v, d2 = c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}
}

func TestDelete(t *testing.T) {
	c := NewCache(CacheOptions{})
	d := c.Set("hello", "world")
	if d != NeverExpires {
		t.Errorf("unexpected deadline: expected NeverExpires, received  %s", d)
	}
	c.Delete("hello")
	v, d2 := c.Get("hello")
	if d2 != d {
		t.Errorf("unexpected deadline: expected %s, received  %s", d, d2)
	}
	if v != nil {
		t.Fatalf("unexpected value: expected <nil>, received %#v", v)
	}
}

func TestPrune(t *testing.T) {
	c := NewCache(CacheOptions{})
	for i := 0; i < 10; i++ {
		c.SetWithTimeout(strconv.Itoa(i), i, time.Millisecond)
	}
	time.Sleep(time.Millisecond)
	c.Prune()
	if c.count() > 0 {
		t.Fatalf("cache should be empty after prune")
	}
}
