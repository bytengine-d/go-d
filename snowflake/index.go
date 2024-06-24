package snowflake

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bytengine-d/go-d/lang"
	"os"
	"strconv"
	"sync"
	"time"
)

const FieldMax = int64(-1) ^ int64(-1)<<int64(5)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1288834974657

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8 = 12

	// DEPRECATED: the below four variables will be removed in a future release.
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

type ID int64

// region Node
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

func (n *Node) Generate() ID {

	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID((now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

func ParseInt64(id int64) ID {
	return ID(id)
}

func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

func ParseIntBytes(id [8]byte) ID {
	return ID(int64(binary.BigEndian.Uint64(id[:])))
}

func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

func (f ID) Step() int64 {
	return int64(f) & stepMask
}

// endregion

// region Global functions
func NewNodeWithGenerateNodeId() (*Node, error) {
	dataCenterId := GetDataCenterId(FieldMax)
	workerId := GetWorkerId(dataCenterId, FieldMax)
	return NewNode(GetNodeId(dataCenterId, workerId))
}

func NewNode(node int64) (*Node, error) {

	if NodeBits+StepBits > 22 {
		return nil, errors.New("Remember, you have a total 22 bits to share between Node/Step")
	}
	// re-calc in case custom NodeBits or StepBits were set
	// DEPRECATED: the below block will be removed in a future release.
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

func GetDataCenterId(max int64) int64 {
	return lang.MacAddrToInt(lang.GetMacAddr(), max)
}

func GetWorkerId(dataCenterId int64, max int64) int64 {
	worker := fmt.Sprintf("%d%d", dataCenterId, os.Getpid())
	var hash int64 = 0
	for _, r := range worker {
		hash = 31*hash + int64(r)
	}
	return (hash & 0xFFFF) % (max + 1)
}

func GetNodeId(dataCenterId, workerId int64) int64 {
	return dataCenterId<<5 | workerId
}

// endregion
