package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/eden-framework/common"
	"github.com/go-redis/redis/v8"
	"github.com/profzone/envconfig"
	"time"
)

type Redis struct {
	Protocol       string
	Host           string
	Port           int
	User           string
	Password       envconfig.Password
	MaxRetries     int
	ConnectTimeout envconfig.Duration
	ReadTimeout    envconfig.Duration
	WriteTimeout   envconfig.Duration
	IdleTimeout    envconfig.Duration
	MinIdle        int
	DB             int

	*redis.Client `ignored:"true"`
}

func (r *Redis) SetDefaults() {
	if r.Protocol == "" {
		r.Protocol = "tcp"
	}
	if r.Port == 0 {
		r.Port = 6379
	}
	if r.MaxRetries == 0 {
		r.MaxRetries = 3
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
		r.IdleTimeout = envconfig.Duration(5 * time.Minute)
	}
	if r.DB == 0 {
		r.DB = 10
	}
}

func (r *Redis) Init() {
	if r.Client != nil {
		return
	}

	r.SetDefaults()
	r.Client = redis.NewClient(&redis.Options{
		Network:      r.Protocol,
		Addr:         fmt.Sprintf("%s:%d", r.Host, r.Port),
		Username:     r.User,
		Password:     r.Password.String(),
		DB:           r.DB,
		MaxRetries:   r.MaxRetries,
		DialTimeout:  time.Duration(r.ConnectTimeout),
		ReadTimeout:  time.Duration(r.ReadTimeout),
		WriteTimeout: time.Duration(r.WriteTimeout),
		MinIdleConns: r.MinIdle,
		IdleTimeout:  time.Duration(r.IdleTimeout),
		TLSConfig:    nil,
		Limiter:      nil,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := r.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis Ping err:%v", err))
	}
}

func (r *Redis) Prefix(key string) string {
	return fmt.Sprintf("%s:%s", "prefix", key)
}

func (r *Redis) Consume(ctx context.Context, topic string, handler func(m common.QueueMessage) error) error {
	if topic == "" {
		return errors.New("[Redis] cannot use Redis as a queue when topic is not specified")
	}
Run:
	for {
		select {
		case <-ctx.Done():
			break Run
		default:
			cmd := r.Client.BRPop(ctx, 0, topic)
			if cmd.Err() != nil {
				return cmd.Err()
			}
			result := cmd.Val()
			m := common.QueueMessage{}
			err := m.UnmarshalBinary([]byte(result[1]))
			if err != nil {
				continue
			}
			err = handler(m)
			if err != nil {
				_ = r.Produce(ctx, m)
			}
		}
	}
	return nil
}

func (r *Redis) Produce(ctx context.Context, messages ...common.QueueMessage) error {
	for _, m := range messages {
		if m.Topic == "" {
			return fmt.Errorf("[Redis] cannot use Redis as a queue when Topic of message is not specified")
		}
		m.Time = time.Now()
	}
	for _, m := range messages {
		cmd := r.Client.LPush(ctx, m.Topic, m)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}
