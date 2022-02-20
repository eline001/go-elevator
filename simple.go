package main

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"test/elevator/conf"
	"test/elevator/elevExec"
	"test/elevator/tool/tool_redis"
	"time"
)

var (
	m sync.Mutex
)

func init() {
	redisErr := tool_redis.NewRedis()
	if redisErr != nil {
		panic(redisErr)
	}
}

var initInputTip = "请输入你想要去的楼层$"

func main() {
	// 初始化程序
	err := elevExec.SimpleExecInit()
	if err != nil {
		fmt.Printf("%v", err)
		panic(err)
	}
	fmt.Printf("程序初始化成功，电梯名称:%s, 最高成为%d 层, 最低层为 %d层 ;"+
		"电梯状态为: %s, 电梯所在楼层为 %d\n", conf.ElevatorShowName, conf.HighestFloorNum, conf.LowestFloorNum, conf.StopIng, conf.CurrentFloorNum)
	// 电梯执行
	go DoWorking(tool_redis.GlobClient)

	// 获取用户在电梯里发出的指令楼层(这里放到一个子线程去完成)
	reader := bufio.NewReader(os.Stdin)
	for {
		//判断用户是否可以输入(是否已经在电梯里了)
		fmt.Print(initInputTip)
		initInputTip = ""
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			panic(err)
		}
		cmdString = strings.Replace(cmdString, "\n", "", -1)
		cmdString = strings.Replace(cmdString, "\r", "", -1)

		// 转换成数字
		var cmdStringInt int
		if cmdStringInt, err = strconv.Atoi(cmdString); err != nil {
			fmt.Println("请输入数字", cmdString, err)
			continue
		} else if cmdStringInt > conf.HighestFloorNum || cmdStringInt < conf.LowestFloorNum {
			fmt.Printf("请输入介于%d到%d之间的楼层\n", conf.LowestFloorNum, conf.HighestFloorNum)
			continue
		}
		userAddFloor(tool_redis.GlobClient, cmdStringInt)
	}
}

// 用户输入楼层信息
func userAddFloor(redisClient *redis.Client, destFloor int) {
	// 获取运行方向
	m.Lock()
	defer m.Unlock()
	err, moveResult := tool_redis.GetMovementDirection(redisClient, conf.ElevatorGroupName, conf.ElevatorUUid)
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return
	}

	// 所在楼层
	err, FloorResult := tool_redis.GetElevatorFloorCache(redisClient, conf.ElevatorGroupName, conf.ElevatorUUid)
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return
	}

	FloorResultInt, err := strconv.Atoi(FloorResult)
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return
	}
	if FloorResultInt == destFloor {
		fmt.Printf("\n不能再当前楼层按当前楼层\n请重新输入$")
		return
	}
	switch moveResult {
	case conf.MoveUp:
		if destFloor-2 < FloorResultInt {
			fmt.Printf("当前运行方向【%s】,所在楼层为 [%d] 新楼层为【%d】,不符合要求，任务丢弃", moveResult, FloorResultInt, destFloor)
			return
		}
		conf.SelfWorkingPool = append(conf.SelfWorkingPool, destFloor)
	case conf.MoveDownward:
		if destFloor > FloorResultInt-2 {
			fmt.Printf("当前运行方向【%s】,所在楼层为 [%d] 新楼层为【%d】,不符合要求，任务丢弃", moveResult, FloorResultInt, destFloor)
			return
		}
		conf.SelfWorkingPool = append(conf.SelfWorkingPool, destFloor)
	case conf.StopIng:
		// 如果是空闲状态，则重新设置方向
		redirect := conf.StopIng
		if destFloor > FloorResultInt {
			// 向上
			redirect = conf.MoveUp
		} else if destFloor < FloorResultInt {
			redirect = conf.MoveDownward
		}
		err = tool_redis.SetMovementDirection(redisClient, conf.ElevatorGroupName, conf.ElevatorUUid, redirect)
		if err != nil {
			fmt.Println("tool_redis.SetMovementDirection", err)
			return
		}
		conf.SelfWorkingPool = append(conf.SelfWorkingPool, destFloor)
	}
	fmt.Printf("任务添加成功%d\n", destFloor)
	return
}

/**
执行任务 (边获取边执行：如果是内部发出的指令，则必须执行完)

内部任务 和 全局任务

优先级：自身的任务 > 全局的任务

1.上行 只获取 上行的任务
2.下行 只获取 下行的任务
3.空闲 获取全部任务
*/
func DoWorking(client *redis.Client) {
	for {
		time.Sleep(time.Second * 2)
		// 获取运行方向
		err, moveResult := tool_redis.GetMovementDirection(client, conf.ElevatorGroupName, conf.ElevatorUUid)
		if err != nil {
			fmt.Println("获取缓存错误moveResult", err)
			panic(err)
		}
		// 所在楼层
		err, FloorResult := tool_redis.GetElevatorFloorCache(client, conf.ElevatorGroupName, conf.ElevatorUUid)
		if err != nil {
			fmt.Println("获取缓存错误FloorResult", err)
			panic(err)
		}
		FloorResultInt, err := strconv.Atoi(FloorResult)
		if err != nil {
			fmt.Println("数据转换错误", err)
			panic(err)
		}
		// 对数据重新排序
		conf.SelfWorkingPool = RemoveRepByLoop(conf.SelfWorkingPool)
		if moveResult == conf.MoveUp {
			sort.Ints(conf.SelfWorkingPool)
		} else if moveResult == conf.MoveDownward {
			sort.Sort(sort.Reverse(sort.IntSlice(conf.SelfWorkingPool)))
		}
		if len(conf.SelfWorkingPool) != 0 {
			fmt.Printf("\n需要到达的楼层 %+v, 当前电梯方向【%s】, 所在楼层【%s】, 行驶速度: XXX,\n ", conf.SelfWorkingPool, moveResult, FloorResult)
			// 判断是否到楼层了----如果到楼成了，则更新数据
			if FloorResult == strconv.Itoa(conf.SelfWorkingPool[0]) {
				fmt.Printf("%s楼到了， %d秒后关门\n", FloorResult, conf.ArrivedStopTime)
				time.Sleep(time.Second * time.Duration(conf.ArrivedStopTime))
				// 删除第一个元素
				conf.SelfWorkingPool = append(conf.SelfWorkingPool[:0], conf.SelfWorkingPool[1:]...)
			}
			// 获取用户在电梯里发出的指令楼层(这里放到一个子线程去完成)
			// 获取相同方向的全局任务（这里放到一个子线程去完成）
			// 电梯执行到下一步
			// 如果到达之后，后续还有楼层则继续下一步
			if len(conf.SelfWorkingPool) != 0 {
				NextFloor := getNextFloor(FloorResultInt, moveResult)
				err = tool_redis.SetElevatorFloorCache(client, conf.ElevatorGroupName, conf.ElevatorUUid, NextFloor)
				if err != nil {
					panic(err)
				}
			}
		} else {
			if moveResult != conf.StopIng {
				err = tool_redis.SetMovementDirection(client, conf.ElevatorGroupName, conf.ElevatorUUid, conf.StopIng)
				if err != nil {
					fmt.Println("tool_redis.SetMovementDirectiontool_redis.SetMovementDirection", err)
					panic(err)
				}
				fmt.Printf("\n 当前停在 %d,请输入你需要到达的楼层$", FloorResultInt)
				continue
			}
		}
	}
}

// 获取下一步的楼层
func getNextFloor(currentFloor int, redirect string) (result int) {
	result = currentFloor
	switch redirect {
	case conf.MoveUp:
		result = 1 + currentFloor
	case conf.MoveDownward:
		result = currentFloor - 1
	}
	//fmt.Printf("resultresultresult %d\n",result)
	return
}

//RemoveRepByLoop  通过两重循环过滤重复元素
func RemoveRepByLoop(slc []int) []int {
	var result []int // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
