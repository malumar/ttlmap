package ttlmap

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// TTL 3 seconds
	seconds := 3
	items := New[int, string](10, seconds, nil)

	for i := 0; i < 10; i++ {
		items.Put(i, fmt.Sprintf("key-%d", i))
	}

	for i := 0; i < 10; i++ {
		if v, ok := items.Get(i); !ok {
			t.Errorf("not found key %v", i)
		} else {
			expected := fmt.Sprintf("key-%d", i)
			if v != expected {
				t.Errorf("got %v expected %v", v, expected)
			}
		}
	}

	time.Sleep(time.Duration(seconds+1) * time.Second)
	for i := 0; i < 10; i++ {
		if v, ok := items.Get(i); ok {
			t.Errorf("found key %v = %v", i, v)
		}
	}

}

func TestNewWithCloser(t *testing.T) {
	// TTL 3 seconds
	seconds := 3
	items := New[int, []byte](10, seconds, func(val *item[[]byte]) {
		fmt.Printf("cleaning %v\n", string(val.Value))
		val.Value = nil
	})

	for i := 0; i < 10; i++ {
		items.Put(i, []byte(fmt.Sprintf("key-%d", i)))
	}

	for i := 0; i < 10; i++ {
		if v, ok := items.Get(i); !ok {
			t.Errorf("not found key %v", i)
		} else {
			expected := []byte(fmt.Sprintf("key-%d", i))

			if bytes.Compare(v, expected) != 0 {
				t.Errorf("got %v expected %v", string(v), string(expected))
			}
		}
	}

	time.Sleep(time.Duration(seconds+1) * time.Second)
	for i := 0; i < 10; i++ {
		if v, ok := items.Get(i); ok {
			t.Errorf("found key %v = %v", i, v)
		}
	}

}
