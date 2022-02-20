package conf

import "fmt"


// 运行方向
const (
	MoveUp       = "上行"       // 向上
	MoveDownward = "下行" // 下行
	StopIng      = "停止状态"      // 停止状态
)



// MovementDirection 运行方向  电梯组+电梯md5id
const MovementDirection = "MovementDirection_%s_%s"

// ElevatorStateCacheKey 电梯状态建 电梯组_电梯唯一标识
const ElevatorStateCacheKey = "state_%s_%s"

//ElevatorFloorCacheKey  电梯所在的层数
const ElevatorFloorCacheKey = "floor_%s_%s"

//ElevatorGroupCacheKey  电梯组包含的电梯
const ElevatorGroupCacheKey = "ElevatorGroupCacheKey_%s"

// ElevatorGlobalWorkPoolUp 所有上行任务池建
const ElevatorGlobalWorkPoolUp = "ElevatorGlobalWorkPoolUp"

// ElevatorGlobalWorkPoolDownward 所有下行任务池建
const ElevatorGlobalWorkPoolDownward = "ElevatorGlobalWorkPoolDownward"

func GetElevatorStateCacheKey(groupName, elevatorName string) string {
	return fmt.Sprintf(ElevatorStateCacheKey, groupName, elevatorName)
}

func GetElevatorFloorCacheKey(groupName, elevatorName string) string {
	return fmt.Sprintf(ElevatorFloorCacheKey, groupName, elevatorName)
}

func GetElevatorGroupCacheKey(groupName string) string {
	return fmt.Sprintf(ElevatorGroupCacheKey, groupName)
}

func GetMovementDirectionKey(groupName,elevatorName string) string {
	
	return fmt.Sprintf(MovementDirection, groupName,elevatorName)
}

