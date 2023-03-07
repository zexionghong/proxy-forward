package gredis

import (
	"encoding/json"
	"fmt"
	"proxy-forward/config"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     config.RuntimeViper.GetInt("redis.max_idle"),
		MaxActive:   config.RuntimeViper.GetInt("redis.max_active"),
		IdleTimeout: time.Duration(config.RuntimeViper.GetInt("redis.idle_timeout")) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.RuntimeViper.GetString("redis.host"), config.RuntimeViper.GetInt("redis.port")), redis.DialPassword(config.RuntimeViper.GetString("redis.password")))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

// Expired a key
func Expired(key string, time int) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	n, _ := redis.Int64(conn.Do("EXPIRE", key, time))
	if n == 1 {
		return true
	}
	return false
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Hmset a key/interface
func Hmset(key string, data map[string]string, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(data)...)

	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

// Hgetall a key
func Hgetall(key string) (map[string]string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Hexists  exists check a hash key
func Hexists(key string, field string) bool {
	conn := RedisConn.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("HEXISTS", redis.Args{}.Add(key).AddFlat(field)...))
	if err != nil {
		return false
	}
	return exists
}

// incr a key/value
func Incr(key string, time int) (int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return 0, err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return 0, err
		}
	}
	return reply, nil
}

// decr  a key/value
func Decr(key string, time int) (int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("DECR", key))
	if err != nil {
		return 0, err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return 0, err
		}
	}
	return reply, nil
}

// incrby
func Incrby(key string, inc int, time int) (int64, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int64(conn.Do("INCRBY", redis.Args{}.Add(key).AddFlat(inc)...))
	if err != nil {
		return 0, err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return 0, err
		}
	}
	return reply, nil
}

/*
@Summary SetMapIncrByInt
@Product SetMapIncr by key
@Params key string
@Params field string
@Params inc int
*/
func SetMapIncrByInt(key string, field string, inc int) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("HINCRBY", redis.Args{}.Add(key).AddFlat(field).AddFlat(inc)...))
	if err != nil {
		return reply, err
	}
	return reply, nil

}

/*
@Summary SetMapIncrByFloat
@Product SetMapIncr by key
@Params key string
@Params field string
@Params inc int
*/
func SetMapIncrByFloat(key string, field string, inc float64) (float64, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Float64(conn.Do("HINCRBYFLOAT", redis.Args{}.Add(key).AddFlat(field).AddFlat(inc)...))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

// Sadd a key value to sadd
func Sadd(key string, value string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("SADD", redis.Args{}.Add(key).AddFlat(value)...))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

/*
@Summary Smembers
@Product return all members in the collection
@Params key string
@return []string
*/
func Smembers(key string) ([]string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

/*
@Summary Lpush
@Product Lpush a valut to list by key
@Params key string
@Params value interface{}
*/
func Lpush(key string, value interface{}) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("LPUSH", redis.Args{}.Add(key).AddFlat(value)...))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

/*
@Summary Llen
@Product Return len(list) by key
@Params key string
@return len
*/
func Llen(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("LLEN", redis.Args{}.Add(key)))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

/*
@Summary Rpop
@Product Return a single of list value
@Params key string
@Return value string
*/
func Rpop(key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("RPOP", redis.Args{}.Add(key)))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

// Zadd a key value to zadd by Int score
func ZaddByInt(key string, score int, value string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("ZADD", redis.Args{}.Add(key).AddFlat(score).AddFlat(value)...))
	if err != nil {
		return reply, err
	}
	return reply, nil
}

// Zadd a key value to zadd by Float score
func ZaddByFloat(key string, score float64, value string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("ZADD", redis.Args{}.Add(key).AddFlat(score).AddFlat(value)...))
	if err != nil {
		return reply, err
	}
	return reply, nil
}
