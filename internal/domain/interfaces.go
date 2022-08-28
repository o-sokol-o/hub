package domain

// Интерфейсы к внешним модулям, которые используются внутри проекта

type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, error)
	Delete(key string)
	Free()
}
