package retrievable

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
) 

type Retrievable interface {
	Key(context.Context, interface{}) *datastore.Key
}

type KeyRetrievable interface {
	Retrievable
	StoreKey(key *datastore.Key)
}

func Serialize(input interface{}) []byte {
	r, _ := json.Marshal(input)
	return r
}

func Unserialize(input []byte, output interface{}) {
	json.Unmarshal(input, output)
}

func PutToMemcache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	mI := &memcache.Item{
		Key:        key,
		Value:      Serialize(value),
		Expiration: expiration,
	}
	return memcache.Set(ctx, mI)
}

func GetFromMemcache(ctx context.Context, output interface{}, key string) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	Unserialize(item.Value, &output)
	return nil
}

func UpdateMemcacheExpire(ctx context.Context, key string, expiration time.Duration) error {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return err
	}
	item.Expiration = expiration
	return memcache.Set(ctx, item)
}

func PlaceInDatastore(ctx context.Context, key interface{}, output Retrievable) (*datastore.Key, error) {
	usrKey := output.Key(ctx, key)
	return datastore.Put(ctx, usrKey, output)
}

func GetFromDatastore(ctx context.Context, key interface{}, output Retrievable) error {
	usrKey := output.Key(ctx, key)
	getErr := datastore.Get(ctx, usrKey, output)
	return getErr
}

func GetData(ctx context.Context, output Retrievable, key interface{}) error {
	DSKey := output.Key(ctx, key)
	MCKey := DSKey.Encode()
	ErrInCache := GetFromMemcache(ctx, output, MCKey)
	if ErrInCache == nil {
		if keyStore, ok := output.(KeyRetrievable); ok {
			keyStore.StoreKey(DSKey)
		}
		return nil
	}
	dataStoreErr := GetFromDatastore(ctx, key, output)
	if dataStoreErr != nil {
		return dataStoreErr
	}
	if keyStore, ok := output.(KeyRetrievable); ok {
		keyStore.StoreKey(DSKey)
	}
	PutToMemcache(ctx, MCKey, output, 0)
	return nil
}

func PlaceData(ctx context.Context, key interface{}, input Retrievable) (*datastore.Key, error) {
	DSKey := input.Key(ctx, key)
	MCKey := DSKey.Encode()
	PutToMemcache(ctx, MCKey, input, 0)
	newKey, dataErr := PlaceInDatastore(ctx, key, input)
	return newKey, dataErr
}

func DeleteEntity(ctx context.Context, key *datastore.Key) error {
	memcache.Delete(ctx, key.Encode())
	return datastore.Delete(ctx, key)
}
