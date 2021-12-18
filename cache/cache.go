package cache

import (
	"fmt"

	"gihtub.com/dshulyak/sqliteexp/database"
	lru "github.com/hashicorp/golang-lru"
)

type conf struct {
	size int
}

type Opt func(*conf)

type Cache struct {
	caches map[string]*lru.Cache
}

func (c *Cache) Register(ns string, opts ...Opt) {
	if _, exist := c.caches[ns]; exist {
		panic(fmt.Sprintf("namespace already exists %s", ns))
	}
	config := &conf{size: 100}
	var err error
	c.caches[ns], err = lru.New(config.size)
	if err != nil {
		panic(err)
	}
}

func (c *Cache) IsCaching(ns string) (database.Accessor, bool) {
	if obj, exist := c.caches[ns]; exist {
		return obj, true
	}
	return nil, false
}
