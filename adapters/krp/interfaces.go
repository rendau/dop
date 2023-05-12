package krp

type Krp interface {
	SendJson(topic, key string, value any) error
	SendManyJson(topic, key string, value []any) error
}
