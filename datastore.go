package retrievable

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Retrievable interface {
	Key(context.Context, interface{}) *datastore.Key
}

type KeyRetrievable interface {
	Retrievable
	StoreKey(key *datastore.Key)
}

func PlaceInDatastore(ctx context.Context, key interface{}, source Retrievable) (*datastore.Key, error) {
	uk := source.Key(ctx, key)
	if uk, putErr := datastore.Put(ctx, uk, source); putErr != nil {
		return uk, putErr
	}
	if i, ok := source.(KeyRetrievable); ok {
		i.StoreKey(uk)
	}
	return uk, putErr
}

func GetFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	if getErr := datastore.Get(ctx, uk, source); getErr != nil {
		return getErr
	}
	if i, ok := source.(KeyRetrievable); ok {
		i.StoreKey(uk)
	}
	return nil
}

func DeleteFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	return datastore.Delete(ctx, uk)
}
