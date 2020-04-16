package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	if err == nil {
		t.Fatalf("should not be connected to redis!")
	}
}

func TestSet(t *testing.T) {
	err := getManager(t).Set("test", "test")
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
	result, err := getManager(t).Get("test")
	require.Equal(t, "test", result, "require test")
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
	result, err := getManager(t).Get("test2")
	require.Equal(t, "test2", result, "require test2")
}

func TestTruncate(t *testing.T) {
	err := getManager(t).Truncate()
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
}

func TestSetWithExpireIfNotExists(t *testing.T) {
	key := "test_ex_nx"
	val := "test"
	err := getManager(t).SetWithExpireIfNotExist(key, val, 10*time.Second)
	if err != nil {
		t.Fatalf("set function doesn't work!")
	}
	result, err := getManager(t).Get(key)
	require.Equal(t, val, result, fmt.Sprintf("require %s", val))

	err = getManager(t).SetWithExpireIfNotExist(key, val, 1*time.Second)
	if err.Error() != "key exists" {
		t.Fatalf("key should exist")
	}
}

func TestGetClient(t *testing.T) {
	require.NotNil(t, getManager(t).GetClient())
}
