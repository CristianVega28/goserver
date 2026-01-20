package it

type StoreInterface interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl int64)
	Delete(key string)
}
