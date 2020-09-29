package redis

import (
	"context"
	"github.com/eden-framework/common"
	"testing"
	"time"
)

func TestRedis_Produce(t *testing.T) {
	cli := &Redis{
		Host:  "localhost",
		Port:  6379,
		Topic: "channel1",
	}
	cli.Init()

	go func() {
		err := cli.Consume(context.Background(), func(m common.QueueMessage) error {
			t.Log(m)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}()
	err := cli.Produce(context.Background(), common.QueueMessage{
		Key:  []byte("foo"),
		Val:  []byte("foo1"),
		Time: time.Now(),
	}, common.QueueMessage{
		Key:  []byte("bar"),
		Val:  []byte("bar1"),
		Time: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {}
}
