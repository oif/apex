package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(512)
	c.Set(123, true, 1)
	v, ok := c.Get(123)
	if !ok {
		t.Fatal("Get cache faile")
	}
	if v != true {
		t.Fatal("Get cache faile")
	}
	if c.Len() != 1 {
		t.Fatal("Wrong cache length")
	}
	time.Sleep(2 * time.Second)
	_, ok = c.Get(123)
	if ok {
		t.Fatal("This cache should expired")
	}

	c.Del(123)
	_, ok = c.Get(123)
	if ok {
		t.Fatal("This cache should be deleted")
	}
	if c.Len() != 0 {
		t.Fatal("Wrong cache length")
	}
}

func BenchmarkWriteCache(b *testing.B) {
	c := New(512)
	for i := 0; i < b.N; i++ {
		c.Set(123, true, 1)
	}
}

func BenchmarkReadCache(b *testing.B) {
	c := New(512)
	c.Set(123, true, 20)

	for i := 0; i < b.N; i++ {
		c.Get(123)
	}
}
