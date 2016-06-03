package retrievable

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"time"
)

var (
	// ErrTooLarge is thrown when this package recognizes that a value for memcache exceeds 1 MB
	ErrTooLarge = errors.New("Esseh/Retrievable Memcache Error: Incoming size is too large, cannot submit to memcache.")
)

func serialize(input interface{}) []byte {
	r, _ := json.Marshal(input)
	return r
}

func unserialize(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
}

// PlaceInMemcache will take a key, value pair and store it in memcache until expiration.
// An error may be returned if the value is too large for memcache to store or if memcache passes an error.
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

// GetFromMemcache will take a key,output struct and attempt to unmarshal the value at key into
// output.
// An error may be returned if json.Unmarshal throws an error or if memcache throws an error.
func GetFromMemcache(ctx context.Context, key string, output interface{}) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	return unserialize(item.Value, &output)
}

// DeleteFromMemcache will attempt to delete memcache memory at key.
// An error may be returned if memcache throws an error
func DeleteFromMemcache(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}

// UpdateMemcacheExpire will attempt to update a value stored at key with new expiration.
// An error may be returned if memcache throws an error
func UpdateMemcacheExpire(ctx context.Context, key string, expiration time.Duration) error {
	item, e := memcache.Get(ctx, key)
	if e != nil {
		return e
	}
	item.Expiration = expiration
	return memcache.Set(ctx, item)
}
