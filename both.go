// Package retrievable handles interaction between
// Google appengine's datastore and memcache using
// a very simple to implement interface.
//
// A trivial example of implementing Retrievable is:
//
//   package main
//
//   import (
//       "github.com/Esseh/retrievable"
//       "golang.org/x/net/context"
//       "google.golang.org/appengine"
//       "google.golang.org/appengine/datastore"
//   )
//
//   type A struct { // A is an example structure implementing KeyRetrievable
//       Value string
//       ID    string `datastore:"-" json:"-"` // We will mark the ID as not to be stored in either Datastore or Memcache
//   }
//
//   func (a *A) Key(ctx context.Context, key interface{}) *datastore.Key {
//       return datastore.NewKey(ctx, "tableA", key.(string), 0, nil) // This is the function to implement Retrievable, A in this context knows the Kind it is a part of and that it's key is a string.
//   }
//
//   func (a *A) StoreKey(key *datastore.Key) {
//       a.ID = key.StringID() // This is the other function to implement KeyRetrievable, now when retrieved A will have it's key in ID.
//   }
//
//   func Example(w http.ResponseWriter, r *http.Request) {
//       ctx := appengine.NewContext(req)
//
//       a := A{}
//       a.Value = "Example Information"
//       // Now that A implements KeyRetrievable, we may use the associated functions.
//       retrievable.PlaceInDatastore(ctx, "Key Value", &a)
//   }
//
// Retrievable structs will follow the properties of datastore.
// If a struct does not implement Serializable, it will also follow the properties of json.
// Information regarding this can be found at:
// https://godoc.org/google.golang.org/appengine/datastore#hdr-Properties
// and
// https://godoc.org/encoding/json#Marshal
//
// More complex and in depth examples can be found in our wiki page at: https://github.com/Esseh/retrievable/wiki
//
package retrievable

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetEntity preforms a get action from first memcache then datastore.
//
// If there is no value in memcache but there is in datastore, this
// function will attempt to place the value in memcache for future
// retrieval.
//
// If found, value is placed into output (Retrievable).
// An error may be returned if datastore cannot find the entity.
func GetEntity(ctx context.Context, key interface{}, output Retrievable) error {
	DSKey := output.Key(ctx, key)

	if cacheErr := GetFromMemcache(ctx, DSKey.Encode(), output); cacheErr == nil {
		if i, ok := output.(KeyRetrievable); ok {
			i.StoreKey(DSKey)
		}
		return nil
	}

	if getErr := GetFromDatastore(ctx, key, output); getErr != nil {
		return getErr
	}

	PlaceInMemcache(ctx, DSKey.Encode(), output, 0)
	return nil
}

// PlaceEntity will place the input Retrievable into first memcache and
// then datastore.
//
// Returns a datastore.Key on successful placement.
// An error may be returned if datastore throws an error.
func PlaceEntity(ctx context.Context, key interface{}, input Retrievable) (*datastore.Key, error) {
	mck := input.Key(ctx, key).Encode()
	PlaceInMemcache(ctx, mck, input, 0)
	return PlaceInDatastore(ctx, key, input)
}

// DeleteEntity will attempt to delete the memory pointed to by key
// first from memcache then datastore.
//
// An error may be returned if datastore throws an error.
func DeleteEntity(ctx context.Context, key *datastore.Key) error {
	DeleteFromMemcache(ctx, key.Encode())
	return datastore.Delete(ctx, key)
}
