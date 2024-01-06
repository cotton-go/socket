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

// Generate 生成一个新的雪花算法 ID
func (s *Snowflake) Generate() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTimestamp := time.Now().UnixNano() / 1e6

	if currentTimestamp < s.lastTimestamp {
		return 0, fmt.Errorf("Clock moved backwards. Refusing to generate ID for %d milliseconds", s.lastTimestamp-currentTimestamp)
	}

	if currentTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence

		if s.sequence == 0 {
			// 等待到下一毫秒
			for currentTimestamp <= s.lastTimestamp {
				currentTimestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = currentTimestamp

	id := (currentTimestamp-epoch)<<(workerIDBits+sequenceBits) | (s.workerID << sequenceBits) | s.sequence
	return id, nil
}
