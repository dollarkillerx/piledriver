package storage

import (
	"github.com/bluele/gcache"
)

var Storage gcache.Cache

func init() {
	gc := gcache.New(20).
		LRU().
		Build()
	Storage = gc
}
