package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	sequence      int64
}

const (
	epoch        = 1609459200000 // 2021-01-01 00:00:00 UTC
	workerIDBits = 5
	sequenceBits = 12
	maxWorkerID  = -1 ^ (-1 << workerIDBits)
	maxSequence  = -1 ^ (-1 << sequenceBits)
)

// NewSnowflake 创建一个新的 Snowflake 实例
//
// 参数：
//   - workerID int64 工作节点的 ID
//
// 返回值：
//   - *Snowflake Snowflake 实例
//   - error 返回错误信息，如果输入的工作节点 ID 在有效范围内
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("Worker ID must be between 0 and %d", maxWorkerID)
	}

	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		sequence:      0,
	}, nil
}

// Generate 生成ID
//
// 参数：
//   - s *Snowflake 雪花算法实例
//
// 返回值：
//   - int64 ID 生成的ID
//   - error 错误信息，如果没有错误则返回nil
func (s *Snowflake) Generate() (int64, error) {
	// 加锁以保证线程安全
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取当前时间戳(毫秒级)
	currentTimestamp := time.Now().UnixNano() / 1e6

	// 如果当前时间戳小于上次生成ID的时间戳，说明时钟回拨了，拒绝生成ID
	if currentTimestamp < s.lastTimestamp {
		return 0, fmt.Errorf("Clock moved backwards. Refusing to generate ID for %d milliseconds", s.lastTimestamp-currentTimestamp)
	}

	// 如果当前时间戳等于上次生成ID的时间戳，序列号加1,并进行位运算处理
	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence

		// 如果序列号为0,等待到下一毫秒再生成ID
		if s.sequence == 0 {
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间戳大于上次生成ID的时间戳，序列号重置为0
		s.sequence = 0
	}

	// 更新上次生成ID的时间戳
	s.lastTimestamp = currentTimestamp

	// 根据公式计算ID,并返回ID和错误信息(如果有的话)
	id := (currentTimestamp-epoch)<<(workerIDBits+sequenceBits) | (s.workerID << sequenceBits) | s.sequence
	return id, nil
}
