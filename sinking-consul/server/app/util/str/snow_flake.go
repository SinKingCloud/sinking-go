package str

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10                      // 每台机器(节点)的ID位数 10位最大可以有2^10=1024个节点
	numberBits  uint8 = 12                      // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位数 即每毫秒可生成 2^12-1=4096个唯一ID
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	numberMax   int64 = -1 ^ (-1 << numberBits) // 同上，用来表示生成id序号的最大值
	timeShift         = workerBits + numberBits // 时间戳向左的偏移量
	workerShift       = numberBits              // 节点ID向左的偏移量
	epoch       int64 = 966441600000
)

var instance *Worker

// GetSnowWorkIns 获取静态对象
func GetSnowWorkIns() *Worker {
	if instance == nil {
		instance, _ = NewSnowWorker(1)
	}
	return instance
}

// Worker 定义一个worker工作节点所需要的基本参数
type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

// NewSnowWorker 实例化一个工作节点
func NewSnowWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker id excess of quantity")
	}
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

// GetId 获取唯一id
func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := (now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID
}

// md5 获取md5值
func (w *Worker) md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUuid 获取uuid
func (w *Worker) GetUuid() string {
	id := w.md5(strconv.FormatInt(w.GetId(), 10))
	id = fmt.Sprintf("%s-%s-%s-%s-%s", id[0:6], id[6:10], id[10:14], id[14:20], id[20:32])
	return id
}
