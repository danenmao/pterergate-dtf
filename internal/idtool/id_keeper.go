package idtool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/internal/redistool"
	"github.com/danenmao/pterergate-dtf/internal/routine"
)

const (
	ReallocStep          = 100
	ReallocThreshold     = 20
	ReallocCheckInterval = 5
)

// id range:
// [Start, FormerEnd),[NewStart, End)
type IdKeeper struct {
	KeyName   string
	Lock      sync.Mutex
	Step      uint32
	Count     uint32
	Start     uint64
	FormerEnd uint64
	NewStart  uint64
	End       uint64
}

var gs_IdKeeper = IdKeeper{}

func GetIdKeeper() *IdKeeper {
	return &gs_IdKeeper
}

func Init(keyName string) error {
	return GetIdKeeper().Init(keyName)
}

func GetId(keyName string) (uint64, error) {
	return GetIdKeeper().GetId(keyName)
}

func (keeper *IdKeeper) Init(keyName string) error {
	if len(keyName) <= 0 {
		return errors.New("empty key name")
	}

	keeper.KeyName = keyName
	keeper.Step = ReallocStep

	// TODO: sync data to MySQL

	// start to maintain the id range
	go func() {
		keeper.maintain()
	}()

	return nil
}

func (keeper *IdKeeper) GetId(keyName string) (uint64, error) {
	if keyName != keeper.KeyName {
		return 0, errors.New("mismatched key name")
	}

	keeper.Lock.Lock()
	defer keeper.Lock.Unlock()

	// check if there is an available id
	for i := 0; i < 2; i++ {
		if keeper.Count <= 0 {
			keeper.Lock.Unlock()

			glog.Info("realloc the id range immediately")
			keeper.realloc()

			keeper.Lock.Lock()
			continue
		} else {
			break
		}
	}

	if keeper.Count <= 0 {
		return 0, errors.New("no available id")
	}

	retId := keeper.updateRange()
	return retId, nil
}

func (keeper *IdKeeper) updateRange() uint64 {
	var retId uint64 = 0
	keeper.Count--

	if keeper.FormerEnd == 0 {
		// a single range, [Start, End)
		retId = keeper.Start
		keeper.Start++

	} else if keeper.FormerEnd != 0 && keeper.Start < keeper.FormerEnd {
		// in a dual range, [Start, FormerEnd)
		retId = keeper.Start
		keeper.Start++

	} else if keeper.FormerEnd != 0 && keeper.Start >= keeper.FormerEnd {
		// switch to [NewStart, End),  move the start pointer to[NewStart, End)
		keeper.Start = keeper.NewStart
		glog.Info("switched to id range: ", keeper.NewStart, keeper.End)

		retId = keeper.Start
		keeper.Start++

		// delete the empty range [Start, FormerEnd)
		keeper.FormerEnd = 0
		keeper.NewStart = 0
	}

	glog.Info(fmt.Sprintf("id status: %d, [%d,%d), %d", retId,
		keeper.Start, keeper.End, keeper.Count))
	return retId
}

func (keeper *IdKeeper) maintain() {
	routine.ExecRoutineWithInterval("realloc",
		func() {
			keeper.realloc()
		},
		time.Second*time.Duration(ReallocCheckInterval))
}

// reallocate the id range
func (keeper *IdKeeper) realloc() {
	if len(keeper.KeyName) <= 0 {
		glog.Warning("empty key name")
		return
	}

	keeper.Lock.Lock()
	defer keeper.Lock.Unlock()

	if keeper.Count > ReallocThreshold {
		return
	}

	// extend the id range
	cmd := redistool.DefaultRedis().IncrBy(context.Background(), keeper.KeyName, ReallocStep)
	val, err := cmd.Result()
	if err != nil {
		glog.Warning("failed to incby ID step: ", err.Error())
		return
	}

	glog.Info("redis incrby id return: ", val)

	keeper.Count += uint32(keeper.Step)

	// [Start, FormerEnd), [NewStart, End)
	keeper.FormerEnd = keeper.End
	keeper.End = uint64(val) + 1

	if keeper.Start == 0 {
		keeper.Start = keeper.End - uint64(keeper.Step)
	}

	keeper.NewStart = keeper.End - uint64(keeper.Step)
	if keeper.FormerEnd > 0 && keeper.FormerEnd-keeper.Start > uint64(keeper.Step) {
		glog.Error("error in id range: [", keeper.Start, keeper.FormerEnd, ")")
	}

	glog.Info(fmt.Sprintf("realloc id range: [%d,%d), [%d,%d), %d",
		keeper.Start, keeper.FormerEnd,
		keeper.NewStart, keeper.End, keeper.Count))
}
