package retrievable

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var (
	// ErrGetZero is thrown when this package sees an invalid get key.
	ErrGetZero = errors.New("Esseh/Retrievable Datastore Error: Zero is not a valid key.")
)

// Retrievable marks a structure as storable with datastore.
type Retrievable interface {

	// Key handles the creation of a datastore key based on
	// an appengine context and unexported information in interface{}.
	//
	// This method must tie to a pointer of your structure. If you are
	// receiving datastore.ErrInvalidEntityType this is the likely issue.
	Key(context.Context, interface{}) *datastore.Key
}

// KeyRetrievable marks a structure as both Retrievable and as storing it's own key.
type KeyRetrievable interface {
	Retrievable

	// StoreKey is a method for the struct to assign it's own key to an internal memory location.
	StoreKey(key *datastore.Key)
}

// PlaceInDatastore will take a Retrievable source and store it into
// datastore based on an appengine context and key information
//
// An error may be returned if datastore passes an error.
// A datastore.Key is returned on a successful put
func PlaceInDatastore(ctx context.Context, key interface{}, source Retrievable) (*datastore.Key, error) {
	uk := source.Key(ctx, key)
	uk, putErr := datastore.Put(ctx, uk, source)
	if putErr != nil {
		return uk, putErr
	}
	if i, ok := source.(KeyRetrievable); ok {
		i.StoreKey(uk)
	}
	return uk, nil
}

// GetFromDatastore will take a Retrievable source and, if possible,
// return the saved struct from datastore.
//
// An error may be returned if datastore passes an error.
//
// Get int64(0) returns the previously accessed object. This is unintended behavior so
// it will also return an error.
// This behavior is insecure, however if the behavior must be used then one level of 
// misdirection will let it through. ie:
//   func (a *A) Key(ctx context.Context, key interface{}) *datastore.Key {
//    return datastore.NewKey(ctx, "tableA", key.(struct{ k int64 }).k, 0, nil)
//   }
	
func GetFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	
	switch key.(type){
		case int64:
			if key.(int64) == 0 {
				return ErrGetZero
			}
	}
	
	if getErr := datastore.Get(ctx, uk, source); getErr != nil {
		return getErr
	}
	if i, ok := source.(KeyRetrievable); ok {
		i.StoreKey(uk)
	}
	return nil
}

// DeleteFromDatastore will take a Retrievable source and, if possible,
// delete the saved struct from datastore.
//
// An error may be returned if datastore passes an error.
func DeleteFromDatastore(ctx context.Context, key interface{}, source Retrievable) error {
	uk := source.Key(ctx, key)
	return datastore.Delete(ctx, uk)
}
