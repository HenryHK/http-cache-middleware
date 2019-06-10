package cache

import (
	"bytes"
	"fmt"
	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// adapterMock is the redis adapter data structure.
type adapterMock struct {
	store *redisCache.Codec
}

// RingOptions exports go-redis RingOptions type.
type RingOptions redis.RingOptions

// Get implements the cache Adapter interface Get method.
func (a *adapterMock) Get(key string) ([]byte, bool) {
	var c []byte
	if err := a.store.Get(key, &c); err != nil {
		return nil, false
	}

	return c, true
}

// Set implements the cache Adapter interface Set method.
func (a *adapterMock) Set(key string, response []byte, expiration time.Time) {
	a.store.Set(&redisCache.Item{
		Key:        key,
		Object:     response,
		Expiration: expiration.Sub(time.Now()),
	})
}

// Release implements the cache Adapter interface Release method.
func (a *adapterMock) Release(key string) {
	fmt.Printf("Delete key: %v\n", key)
	err := a.store.Delete(key)
	if err != nil {
		fmt.Printf("error deleting redis: %v\n", err)
	}
}

// New initializes Redis adapter.
func newMock(opt *RingOptions) *adapterMock {
	ropt := redis.RingOptions(*opt)
	return &adapterMock{
		&redisCache.Codec{
			Redis: redis.NewRing(&ropt),
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)

			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		},
	}
}

func TestMiddleware(t *testing.T) {

	counter := 0
	httpTestHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("new value %v", counter)))
	})

	opt := &RingOptions{
		Addrs: map[string]string{
			"server": "localhost:6379",
		},
	}

	client, _ := New(
		ClientWithAdapter(newMock(opt)),
		ClientWithTTL(10*time.Minute),
	)

	handler := client.Middleware(httpTestHandler)

	tests := []struct {
		name     string
		url      string
		method   string
		body     []byte
		wantBody string
		wantCode int
	}{
		{
			"returns cached response",
			"http://foo.bar/contact/test-1",
			"GET",
			nil,
			"new value 1",
			200,
		},
		{
			"release cache and update",
			"http://foo.bar/contact",
			"POST",
			[]byte("{\n  \"contact\": {\n    \"FirstName\": \"Slarty\",\n    \"LastName\": \"Bartfast\",\n    \"Email\": \"test@slarty.com\",\n    \"custom\": {\n      \"string--Test--Field\": \"This is a test\"\n    }\n  }\n}"),
			"new value 2",
			200,
		},
		{
			"release cache and update via POST",
			"http://foo.bar/contact",
			"PUT",
			[]byte("{\n  \"contact\": {\n    \"FirstName\": \"Slarty\",\n    \"LastName\": \"Bartfast\",\n    \"Email\": \"test@slarty.com\",\n    \"custom\": {\n      \"string--Test--Field\": \"This is a test\"\n    }\n  }\n}"),
			"new value 3",
			200,
		},
		{
			"release cache",
			"http://foo.bar/contact/test-3",
			"DELETE",
			nil,
			"new value 4",
			200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter++

			r, err := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
			if err != nil {
				t.Error(err)
				return
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if !reflect.DeepEqual(w.Code, tt.wantCode) {
				t.Errorf("*Client.Middleware() = %v, want %v", w.Code, tt.wantCode)
				return
			}
			if !reflect.DeepEqual(w.Body.String(), tt.wantBody) {
				t.Errorf("*Client.Middleware() = %v, want %v", w.Body.String(), tt.wantBody)
			}
		})
	}

}
