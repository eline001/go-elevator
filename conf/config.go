package conf

// ElevatorUUid 电梯名称
var ElevatorUUid string

// ElevatorShowName 电梯显示名称
var ElevatorShowName string

//ElevatorGroupName  电梯组名称
var ElevatorGroupName string

// HighestFloorNum 楼层最高层
var HighestFloorNum int

// LowestFloorNum 楼层最底层
var LowestFloorNum int

// CurrentFloorNum 当前电梯所在楼层
var CurrentFloorNum int

// IsAtElevator 是否在电梯里
var IsAtElevator bool

// SelfWorkingPool 电梯自身的任务池
var SelfWorkingPool []int

//  ArrivedStopTime 电梯到了，几秒后关门
var ArrivedStopTime = 5

// 运动反方向
var MoveReverse = map[string]string{
	MoveUp:       MoveDownward, // 向上
	MoveDownward: MoveUp,       // 下行
	StopIng:      StopIng,
}


// 本地维护当前电梯状态
var CurrentElevatorState = ""