package redishelper

type Redis interface {
	connectRedis()

	Read(typeDB int, key string) (any, error)
	CreateUpdate(typeDB int, key string, value any) (string, error)
	Delete(typeDB int, key string) error
}
