package utils

import (
	"sync"
	"time"
)

// 支持 2 ^ 8 - 1 台机器
//每一个毫秒支持 2 ^ 9 - 1 个不同的id
const (
	workerBits   = uint(10)
	maxWorkerId  = int64(-1 ^ (-1 << workerBits))
	StepBits     = uint(12)
	maskStep     = int64(-1 ^ (-1 << StepBits))
	maskWorkerId = workerBits << StepBits

	timeShift   = workerBits + StepBits
	workerShift = StepBits
)

// Worker 定义一个woker工作节点所需要的基本参数
type Worker struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	workerId  int64      // 机器编码
	timestamp int64      // 记录时间戳
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	epoch     time.Time  //初始时间
}

// NewSnowflakeWorker NewWorker 初始化ID生成结构体
// workerId 机器的编号
func NewSnowflakeWorker(workerId int64) *Worker {
	if workerId > maxWorkerId {
		panic("workerId 不能大于最大值")
	}

	var curTime = time.Now()
	epoch := int64(1672502400000)

	return &Worker{
		workerId:  workerId,
		timestamp: 0,
		number:    0,
		epoch:     curTime.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(curTime))}
}

// GetId 生成id 的方法用于生成唯一id
func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Since(w.epoch).Milliseconds() // time.Now().UnixMilli() // 获得现在对应的时间戳
	if now < w.timestamp {
		// 当机器出现时钟回拨时报错
		panic("Clock moved backwards.  Refusing to generate id for %d milliseconds")
	}
	if w.timestamp == now {
		w.number = (w.number + 1) & maskStep

		if w.number == 0 { //此处为最大节点ID,大概是2^9-1 511条,
			for now <= w.timestamp {
				now = time.Since(w.epoch).Milliseconds()
			}
		}
	} else {
		w.number = 0
	}
	w.timestamp = now // 将机器上一次生成ID的时间更新为当前时间
	ID := int64((now)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}
