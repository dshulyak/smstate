package database

func AddIfCaching(db Executor, ns string, key, value interface{}) {
	switch caching := db.(type) {
	case Caching:
		cache, enabled := caching.IsCaching(ns)
		if enabled {
			cache.Add(key, value)
		}
	}
}

type Caching interface {
	IsCaching(string) (Accessor, bool)
}

type Accessor interface {
	Add(interface{}, interface{}) bool
	Get(interface{}) (interface{}, bool)
}
