package adapter

import (
	"fmt"
	cache "github.com/HenryHK/http-cache-middleware/cache"
	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

// Adapter is the memory adapter data structure.
type Adapter struct {
	store *redisCache.Codec
}

// RingOptions exports go-redis RingOptions type.
type RingOptions redis.RingOptions

// Get implements the cache Adapter interface Get method.
func (a *Adapter) Get(key string) ([]byte, bool) {
	var c []byte
	if err := a.store.Get(key, &c); err != nil {
		return nil, false
	}

	return c, true
}

// Set implements the cache Adapter interface Set method.
func (a *Adapter) Set(key string, response []byte, expiration time.Time) {
	a.store.Set(&redisCache.Item{
		Key:        key,
		Object:     response,
		Expiration: expiration.Sub(time.Now()),
	})
}

// Release implements the cache Adapter interface Release method.
func (a *Adapter) Release(key string) {
	fmt.Printf("Delete key: %v\n", key)
	err := a.store.Delete(key)
	if err != nil {
		log.Printf("error deleting redis: %v\n", err)
	}
}

// New initializes Redis adapter.
func New(opt *RingOptions) cache.Adapter {
	ropt := redis.RingOptions(*opt)
	return &Adapter{
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
