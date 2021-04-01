package rediser

import (
	"sync"

	"github.com/go-redis/redis"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/basic/config"
)

var (
	client *redis.Client
	m      sync.RWMutex
	inited bool
)

// Init make the redis client avaiable
func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Info("Already inited Redis...")
		return
	}

	redisConfig := config.GetRedisConfig()

	if redisConfig != nil && redisConfig.GetEnabled() {
		log.Info("Init Redis...")

		if redisConfig.GetSentinelConfig() != nil && redisConfig.GetSentinelConfig().GetEnabled() {
			log.Info("Init Redis in SentinelMode...")
			initSentinel(redisConfig)
		} else {
			log.Info("Init Redis in NormalMode...")
			initSingle(redisConfig)
		}

		log.Info("Redis checking connection...")

		pong, err := client.Ping().Result()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Info("Redis Ping.")
		log.Info("Redis Ping..")
		log.Info("Redis Ping...", pong)
	}
	inited = true
	log.Info("Redis connected successfully!")
}

func GetRedis() *redis.Client {
	return client
}

func initSentinel(redisConfig config.RedisConfig) {
	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    redisConfig.GetSentinelConfig().GetMaster(),
		SentinelAddrs: redisConfig.GetSentinelConfig().GetNodes(),
		DB:            redisConfig.GetDBNum(),
		Password:      redisConfig.GetPassword(),
	})

}

func initSingle(redisConfig config.RedisConfig) {
	client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.GetConn(),
		Password: redisConfig.GetPassword(), // no password set
		DB:       redisConfig.GetDBNum(),    // use default DB
	})
}
