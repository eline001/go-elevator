package tool_redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"test/elevator/conf"
)

const (
	KeyNotExist = "KeyNotExist"
)

type RedisServiceErr struct {
	msg string
}

func (rse RedisServiceErr) Error() string {
	return rse.msg
}

func NewRedisServiceErr(msg string) error {
	return &RedisServiceErr{
		msg: msg,
	}
}

func GetElevatorStateCache(client *redis.Client, groupName, ElevatorUUid string) (err error, str string) {
	floorResult := client.Get(context.Background(), conf.GetElevatorStateCacheKey(groupName, ElevatorUUid))
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// GetElevatorFloorCache 电梯所在楼层
func GetElevatorFloorCache(client *redis.Client, groupName, ElevatorUUid string) (err error, str string) {
	floorResult := client.Get(context.Background(), conf.GetElevatorFloorCacheKey(groupName, ElevatorUUid))
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// SetElevatorFloorCache 设置电梯所在楼层
func SetElevatorFloorCache(client *redis.Client, groupName, ElevatorUUid string, floor int) error {
	floorResult := client.Set(context.Background(), conf.GetElevatorFloorCacheKey(groupName, ElevatorUUid), floor, -1)
	return floorResult.Err()
}

func GetElevatorGroupCache(client *redis.Client, groupName, ElevatorUUid string) (err error, str string) {
	floorResult := client.HGet(context.Background(), conf.GetElevatorGroupCacheKey(groupName), ElevatorUUid)
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// GetElevatorGlobalWorkPoolUp 获取所有上行任务池
func GetElevatorGlobalWorkPoolUp(client *redis.Client) (err error, str string) {
	floorResult := client.Get(context.Background(), conf.ElevatorGlobalWorkPoolUp)
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// GetElevatorGlobalWorkPoolDownward 获取所有下行任务池
func GetElevatorGlobalWorkPoolDownward(client *redis.Client) (err error, str string) {
	floorResult := client.Get(context.Background(), conf.ElevatorGlobalWorkPoolDownward)
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// GetMovementDirection 获取电梯方向
func GetMovementDirection(client *redis.Client, groupName, ElevatorUUid string) (err error, str string) {
	floorResult := client.Get(context.Background(), conf.GetMovementDirectionKey(groupName, ElevatorUUid))
	if floorResult.Err() != nil {
		if !errors.Is(floorResult.Err(), redis.Nil) {
			return NewRedisServiceErr(KeyNotExist), ""
		}
		return floorResult.Err(), ""
	}
	str = floorResult.Val()
	return
}

// SetMovementDirection  设置电梯方向
func SetMovementDirection(client *redis.Client, groupName, ElevatorUUid string, val string) error {
	conf.CurrentElevatorState = val
	floorResult := client.Set(context.Background(), conf.GetMovementDirectionKey(groupName, ElevatorUUid), val, -1)
	return floorResult.Err()
}
