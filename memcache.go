// Package retrievable handles interaction between
// Google appengine's datastore and memchache using
// a very simple to implement interface.
package retrievable

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"time"
)

var (
	ErrTooLarge = errors.New("Esseh/Retrivable Memchache Error: Incoming size is too large, cannot submit to memcache.")
)

func serialize(input interface{}) []byte {
	r, _ := json.Marshal(input)
	return r
}

func unserialize(input []byte, output interface{}) {
	json.Unmarshal(input, output)
}

func PlaceInMemcache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	v := serialize(value)
	if len(v) > 10e5 { // Memcache limits to 1 MB
		return ErrTooLarge
	}
	mI := &memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: expiration,
	}
	return memcache.Set(ctx, mI)
}

func GetFromMemcache(ctx context.Context, key string, output interface{}) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	unserialize(item.Value, &output)
	return nil
}

func DeleteFromMemcache(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}

func UpdateMemcacheExpire(ctx context.Context, key string, expiration time.Duration) error {
	item, e := memcache.Get(ctx, key)
	if e != nil {
		return e
	}
	item.Expiration = expiration
	return memcache.Set(ctx, item)
}
