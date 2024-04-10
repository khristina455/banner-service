package cache

type RedisClient struct {
	//конфиг
	//клиент
}

func NewRedisClient() *RedisClient { //передаем конфиг редиски
	return &RedisClient{}
}

func Get(key string) (value []byte, ok bool) {
	return []byte{}, false
}

func Set(key string, value []byte) {

}
