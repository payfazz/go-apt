package redis

import (
	"testing"
	"time"
)

func getManager(t *testing.T) RedisInterface {
	manager, err := NewFazzRedis("localhost:6379", "cashfazz")
	if err != nil {
		t.Fatalf("cannot connect to redis!")
	}
	return manager
}

func TestFailingConnect(t *testing.T) {
	_, err := NewFazzRedis("localhost:6379", "")
	if err != nil {
		t.Fatalf("cannot connect to redis!")
	}
}

func TestSet(t *testing.T) {
	err := getManager(t).Set("test", "test")
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}

func TestGet(t *testing.T) {
	_, err := getManager(t).Get("test")
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}

func TestDelete(t *testing.T) {
	err := getManager(t).Delete("test")
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}

func TestSetWithExpire(t *testing.T) {
	err := getManager(t).SetWithExpire("test2", "test2", 1*time.Second)
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}

func TestTruncate(t *testing.T) {
	err := getManager(t).Truncate()
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}
