package redix2

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
	"weixingongzhonghao/config"
)

type Redix struct {
	RedisPoolAct *redis.Pool // redis实例
}

var (
	RedisX *Redix
)

func InitRedis() {
	RedisX = &Redix{}
	RedisX.initRedixAct()
}

func (redis *Redix) initRedixAct() {
	//连接redis
	//RedisX = &Redix{
	//	RedisPoolAct: Conn(viper.GetString("redis_act.conn"),viper.GetString("redis_act.passwd"),viper.GetInt("redis_act.dbnum")),
	//}
	redis.RedisPoolAct = Conn(config.G_JsonConfig.Redis.Host, config.G_JsonConfig.Redis.Password, config.G_JsonConfig.Redis.Database)

	r := redis.RedisPoolAct.Get()
	if r.Err() != nil {
		log.Fatal(r.Err())
	}
	defer r.Close()

}

// 获取次
func (this *Redix) GetRedixAct() redis.Conn {
	r := this.RedisPoolAct.Get()
	return r
}

func Conn(conn, auth string, dbnum int) *redis.Pool {
	return &redis.Pool{
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 180 * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				conn,
				redis.DialConnectTimeout(time.Duration(5)*time.Second),
				redis.DialReadTimeout(time.Duration(10)*time.Second),
				redis.DialWriteTimeout(time.Duration(10)*time.Second),
			)
			if err != nil {
				return nil, err
			}
			//验证redis 是否有密码
			passwd := auth
			if passwd != "" {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			c.Do("select", dbnum)

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

}

func (this *Redix) CommonCmdAct(cmdStr string, keysAndArgs ...interface{}) (reply interface{}, err error) {

	c := this.GetRedixAct()
	if c.Err() != nil {
		return
	}
	defer c.Close()
	res, err := c.Do(cmdStr, keysAndArgs...)
	if err == nil {
		return res, err
	}

	return nil, err
}

func (r *Redix) Get(k string) (interface{}, error) {
	c := r.RedisPoolAct.Get()
	defer c.Close()
	v, err := c.Do("GET", k)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *Redix) Set(k, v string) error {
	c := r.RedisPoolAct.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	return err
}

func (r *Redix) SetEx(k string, v interface{}, ex int64) error {
	c := r.RedisPoolAct.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v, "EX", ex)
	return err
}
