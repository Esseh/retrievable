package retrievable

import (
	"google.golang.org/appengine/datastore"
)

// IntID is a shortcut type that can be embeded in another struct to fulfil
// the KeyRetrievable interface easily in the most common case.
type IntID int64

func (i *IntID) StoreKey(key *datastore.Key) {
	*i = IntID(key.IntID())
}

// StringID is a shortcut type that can be embeded in another struct to fulfil
// the KeyRetrievable interface easily in the most common case.
type StringID string

func (s *StringID) StoreKey(key *datastore.Key) {
	*s = StringID(key.StringID())
}
