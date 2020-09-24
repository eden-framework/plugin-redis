package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/profzone/envconfig"
	"time"
)

type Redis struct {
	Protocol       string
	Host           string
	Port           int
	User           string
	Password       envconfig.Password
	ConnectTimeout envconfig.Duration
	ReadTimeout    envconfig.Duration
	WriteTimeout   envconfig.Duration
	IdleTimeout    envconfig.Duration
	MaxActive      int
	MaxIdle        int
	Wait           bool
	DB             int

	client *redis.Pool
}

func (r *Redis) SetDefaults() {
	if r.Protocol == "" {
		r.Protocol = "tcp"
	}
	if r.Port == 0 {
		r.Port = 6379
	}
	if r.ConnectTimeout == 0 {
		r.ConnectTimeout = envconfig.Duration(10 * time.Second)
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = envconfig.Duration(10 * time.Second)
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = envconfig.Duration(10 * time.Second)
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = envconfig.Duration(240 * time.Second)
	}
	if r.MaxActive == 0 {
		r.MaxActive = 5
	}
	if r.MaxIdle == 0 {
		r.MaxIdle = 3
	}
	if !r.Wait {
		r.Wait = true
	}
	if r.DB == 0 {
		r.DB = 10
	}
}

func (r *Redis) Init() {
	if r.client != nil {
		return
	}

	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial(
			r.Protocol,
			fmt.Sprintf("%s:%d", r.Host, r.Port),

			redis.DialWriteTimeout(time.Duration(r.WriteTimeout)),
			redis.DialConnectTimeout(time.Duration(r.ConnectTimeout)),
			redis.DialReadTimeout(time.Duration(r.ReadTimeout)),
			redis.DialUsername(r.User),
			redis.DialPassword(r.Password.String()),
			redis.DialDatabase(r.DB),
		)
		return c, err
	}

	r.client = &redis.Pool{
		Dial:        dialFunc,
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxActive,
		IdleTimeout: time.Duration(r.IdleTimeout),
		Wait:        true,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 5*time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Prefix(key string) string {
	return fmt.Sprintf("%s:%s", "prefix", key)
}

func (r *Redis) Get() redis.Conn {
	if r.client != nil {
		return r.client.Get()
	}
	return nil
}

func (r *Redis) Query(cmd *CMD, others ...*CMD) (interface{}, error) {
	c := r.Get()
	defer c.Close()

	if (len(others)) == 0 {
		return c.Do(cmd.name, cmd.args...)
	}

	err := c.Send("MULTI")
	if err != nil {
		return nil, err
	}

	err = c.Send(cmd.name, cmd.args...)
	if err != nil {
		return nil, err
	}

	for i := range others {
		o := others[i]
		if o == nil {
			continue
		}
		err := c.Send(o.name, o.args...)
		if err != nil {
			return nil, err
		}
	}

	return c.Do("EXEC")
}

func (r *Redis) Execute(cmd *CMD) error {
	c := r.Get()
	defer c.Close()

	return c.Send(cmd.name, cmd.args...)
}
