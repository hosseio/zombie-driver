package zombie_configuration

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/heetch/jose-odg-technical-test/zombie-driver/driver-zombie"
)

const (
	key = "heetch.zombie-configuration"
)

type RedisZombieConfigurer struct {
	client redis.Pool
}

type RedisAddr string

func NewRedisZombieConfigurer(address RedisAddr) RedisZombieConfigurer {
	pool := redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", string(address))
		},
		MaxIdle:   80,
		MaxActive: 12000,
	}

	return RedisZombieConfigurer{pool}
}

func (rc RedisZombieConfigurer) GetZombieConfig() (driver_zombie.ZombieConfigProjection, error) {
	conn := rc.client.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return driver_zombie.ZombieConfigProjection{}, err
	}

	var zombieConfig driver_zombie.ZombieConfigProjection
	err = json.Unmarshal(bytes, &zombieConfig)
	if err != nil {
		return driver_zombie.ZombieConfigProjection{}, err
	}

	return zombieConfig, nil
}

func (rc RedisZombieConfigurer) SetZombieConfig(config driver_zombie.ZombieConfigProjection) error {
	conn := rc.client.Get()
	defer conn.Close()

	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, bytes)

	return err
}
