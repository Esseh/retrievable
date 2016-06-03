// Package retrievable handles interaction between
// Google appengine's datastore and memchache using
// a very simple to implement interface.
package retrievable

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Retrievable marks a structure as interactable with datastore.
type Retrievable interface {

	// Key handles the creation of a datastore key based on
	// an appengine context and unexported information in interface{}.
	//
	// This method must tie to a pointer of your structure. If you are
	// receiving datastore.ErrInvalidEntityType this is the likey issue.
	Key(context.Context, interface{}) *datastore.Key
}

// KeyRetrievable marks a structure as both Retrievable and as storing it's own key.
type KeyRetrievable interface {
	Retrievable

	// StoreKey is a method for the struct to assign it's own key to an internal memory location.
	StoreKey(key *datastore.Key)
}

// PlaceInDatastore will take a Retrievable source and store it into datastore based on an appengine context and key information
// An error may be returned if datastore passes an error.
// The datastore.Key that is returned from a successful push to datatastore is also passed up.
func PlaceInDatastore(ctx context.Context, key interface{}, source Retrievable) (*datastore.Key, error) {
	uk := source.Key(ctx, key)
	if uk, putErr := datastore.Put(ctx, uk, source); putErr != nil {
		return uk, putErr
	}
	if i, ok := source.(KeyRetrievable); ok {
		i.StoreKey(uk)
	}
	return uk, nil
}

// GetFromDatastore will take a Retrievable source and, if possible, return the saved struct from datastore.
// An error may be returned if datastore passes an error.
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

// DeleteFromDatastore
func DeleteFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	return datastore.Delete(ctx, uk)
}
