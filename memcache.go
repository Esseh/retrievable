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

// Serializable marks a structure as having a custom storage implementation in memcache.
//
// If a struct implements Serializable, retrievable will use its methods instead of json for memcache.
type Serializable interface {
	// Serialize should create a slice of bytes that represents the structure.
	Serialize() []byte
	// Unserialize should take the slice of bytes from serialize and recreate the structure.
	Unserialize([]byte) error
}

func serialize(input interface{}) []byte {
	r, _ := json.Marshal(input)
	return r
}

func unserialize(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
}

// PlaceInMemcache will take a key, value pair and store it in memcache
// until expiration.
//
// An error may be returned if the value is too large for memcache to
// store or if memcache passes an error.
func PlaceInMemcache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var v []byte
	if ser, ok := value.(Serializable); ok {
		v = ser.Serialize()
	} else {
		v = serialize(value)
	}
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

// GetFromMemcache will take a key,output struct and attempt to
// unmarshal the value at key into output.
//
// An error may be returned if json.Unmarshal throws an error or
// if memcache throws an error.
func GetFromMemcache(ctx context.Context, key string, output interface{}) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	if ser, ok := output.(Serializable); ok {
		return ser.Unserialize(item.Value)
	}
	return unserialize(item.Value, &output)
}

// DeleteFromMemcache will attempt to delete memcache memory at key.
// An error may be returned if memcache throws an error
func DeleteFromMemcache(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}

// UpdateMemcacheExpire will attempt to update a value stored at
// key with new expiration.
//
// An error may be returned if memcache throws an error
func UpdateMemcacheExpire(ctx context.Context, key string, expiration time.Duration) error {
	item, e := memcache.Get(ctx, key)
	if e != nil {
		return e
	}
	item.Expiration = expiration
	return memcache.Set(ctx, item)
}
