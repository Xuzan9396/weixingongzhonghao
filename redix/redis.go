package redix

import (
	"github.com/silenceper/wechat/v2/cache"
	"weixingongzhonghao/config"
)

type RedisPool struct {
	Redis *cache.Redis
}

var g_redis *RedisPool

func GetReids() *RedisPool {

	if g_redis == nil {
		g_redis = &RedisPool{}
		redisOpts := &cache.RedisOpts{
			Host:        config.G_JsonConfig.Redis.Host,     // redis host
			Password:    config.G_JsonConfig.Redis.Password, //redis password
			Database:    config.G_JsonConfig.Redis.Database, // redis db
			MaxActive:   10,                                 // 连接池最大活跃连接数
			MaxIdle:     10,                                 //连接池最大空闲连接数
			IdleTimeout: 60,                                 //空闲连接超时时间，单位：second
		}
		//fmt.Println(cache.NewRedis(redisOpts) == nil );
		//redisPool.Redis =
		g_redis.Redis = cache.NewRedis(redisOpts)

	}

	return g_redis

}
