// Package retrievable handles interaction between
// Google appengine's datastore and memchache using
// a very simple to implement interface.
//
// More Documentation
package retrievable

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func GetEntity(ctx context.Context, output Retrievable, key interface{}) error {
	DSKey := output.Key(ctx, key)

	if cacheErr := GetFromMemcache(ctx, DSKey.Encode(), output); cacheErr == nil {
		if i, ok := output.(KeyRetrievable); ok {
			i.StoreKey(DSKey)
		}
		return nil
	}

	if getErr := GetFromDatastore(ctx, DSKey, output); getErr != nil {
		return getErr
	}

	PlaceInMemcache(ctx, DSKey.Encode(), output, 0)
	return nil
}

func PlaceEntity(ctx context.Context, key interface{}, input Retrievable) (*datastore.Key, error) {
	mck := input.Key(ctx, key).Encode()
	PlaceInMemcache(ctx, mck, input, 0)
	return PlaceInDatastore(ctx, key, input)
}

func DeleteEntity(ctx context.Context, key *datastore.Key) error {
	DeleteFromMemcache(ctx, key.Encode())
	return datastore.Delete(ctx, key)
}
