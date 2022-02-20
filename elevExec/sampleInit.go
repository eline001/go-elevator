package elevExec

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	mathRand "math/rand"
	"test/elevator/conf"
	"test/elevator/tool/tool_redis"
	"time"
)

/**
2.2 初始化电梯数量
2.3 全局记录电梯所在楼层
2.4 全局任务池(乘客在电梯外发出的指令)
2.5 各个电梯的状态(空闲，工作，抛锚)
*/

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}
func NewError(msg string) error {
	return &Error{msg: msg}
}

// ExecInit
/**
  初始化程序
  1. 初始化电梯数量
  2. 将电梯所在层数全部初始化在任意层
  3. 设置电梯的状态为空闲
*/

//SimpleExecInit 初始化电梯
/**
    1.设置当前电梯状态
    2.设置当前电梯所在层数
	3.将电梯放入电梯组中
*/
func SimpleExecInit() error {
	flag.StringVar(&conf.ElevatorGroupName, "gn", "电梯1组", "电梯组名称")
	flag.StringVar(&conf.ElevatorShowName, "sn", "一号电梯", "显示名称")
	flag.IntVar(&conf.HighestFloorNum, "fn", 40, "最高楼层")
	flag.IntVar(&conf.LowestFloorNum, "lfn", -2, "最低楼层")
	flag.Parse()
	if conf.ElevatorGroupName == "" || conf.HighestFloorNum == 0 || conf.ElevatorShowName == "" {
		flag.PrintDefaults()
		fmt.Println("请按照参数输入数据")
		return NewError("请按照参数输入数据")
	}

	mathRand.Seed(int64(time.Now().Nanosecond()))
	conf.ElevatorUUid = uuid.NewString()

	// 1.设置当前电梯状态
	setResult := tool_redis.GlobClient.Set(context.Background(), conf.GetElevatorStateCacheKey(conf.ElevatorGroupName, conf.ElevatorUUid),
		conf.StopIng, -1)
	if setResult.Err() != nil {
		return setResult.Err()
	}

	// 2.设置当前电梯所在层数
	conf.CurrentFloorNum = mathRand.Intn(conf.HighestFloorNum)
	setResult = tool_redis.GlobClient.Set(context.Background(), conf.GetElevatorFloorCacheKey(conf.ElevatorGroupName, conf.ElevatorUUid),
		conf.CurrentFloorNum, -1)
	if setResult.Err() != nil {
		return setResult.Err()
	}
	// 3.将电梯放入电梯组中
	iResult := tool_redis.GlobClient.HSet(context.Background(), conf.GetElevatorGroupCacheKey(conf.ElevatorGroupName), conf.ElevatorUUid,
		conf.ElevatorShowName)
	if iResult.Err() != nil {
		return iResult.Err()
	}

	// 设置电梯运行方向
	err := tool_redis.SetMovementDirection(tool_redis.GlobClient, conf.ElevatorGroupName, conf.ElevatorUUid, conf.StopIng)
	if err != nil {
		return err
	}
	return nil
}
