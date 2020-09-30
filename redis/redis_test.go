package redis

import (
	"context"
	"github.com/eden-framework/common"
	"testing"
)

func TestRedis_Produce(t *testing.T) {
	cli := &Redis{
		Host: "localhost",
		Port: 6379,
	}
	cli.Init()

	go func() {
		err := cli.Consume(context.Background(), "channel1", func(m common.QueueMessage) error {
			t.Log(m)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}()
	err := cli.Produce(context.Background(), common.QueueMessage{
		Topic: "channel1",
		Key:   []byte("foo"),
		Val:   []byte("foo1"),
	}, common.QueueMessage{
		Topic: "channel1",
		Key:   []byte("bar"),
		Val:   []byte("bar1"),
	})
	if err != nil {
		t.Fatal(err)
	}

	select {}
}

type testStruct string

func (t *testStruct) UnmarshalBinary(data []byte) error {
	*t = testStruct(data)
	return nil
}

func (t testStruct) MarshalBinary() (data []byte, err error) {
	return []byte(t), nil
}

func TestGetAndSet(t *testing.T) {
	cli := &Redis{
		Host: "localhost",
		Port: 6379,
	}
	cli.Init()

	err := cli.Set(context.Background(), "foo", testStruct("bar"), 0)
	if err != nil {
		t.Fatal(err)
	}

	var result testStruct
	err = cli.Get(context.Background(), "foo", &result)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)

	err = cli.Del(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}
}
