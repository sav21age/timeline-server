package cache

import "errors"

var ErrItemNotFound = errors.New("cache item not found")

type Cache interface {
	Set(key, value interface{}, ttl int64) error
	Get(key interface{}) (interface{}, error)
}
