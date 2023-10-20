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

// ID的步长设置
const (
	ReallocStep          = 100
	ReallocThreshold     = 20
	ReallocCheckInterval = 5
)

// ID可用范围：[Start, FormerEnd),[NewStart, End)
type IdKeeper struct {
	KeyName   string
	Lock      sync.RWMutex
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

// 初始化
func Init(keyName string) error {
	return GetIdKeeper().Init(keyName)
}

// 获取可用的ID
func GetId(keyName string) (uint64, error) {
	return GetIdKeeper().GetId(keyName)
}

func (keeper *IdKeeper) Init(keyName string) error {
	if len(keyName) <= 0 {
		return errors.New("empty key name")
	}

	keeper.KeyName = keyName
	keeper.Step = ReallocStep

	// TODO: 定期同步至MySQL, 处理迁移的场景

	// 启动后台维护协程
	go func() {
		keeper.refreshId()
	}()

	return nil
}

func (keeper *IdKeeper) GetId(keyName string) (uint64, error) {
	if keyName != keeper.KeyName {
		return 0, errors.New("mismatched key name")
	}

	keeper.Lock.Lock()
	defer keeper.Lock.Unlock()

	// 检查是否有可用ID
	for i := 0; i < 2; i++ {
		if keeper.Count <= 0 {
			keeper.Lock.Unlock()

			glog.Info("renew id range immediately")
			keeper.reallocIdIfNeed()

			keeper.Lock.Lock()
			continue
		} else {
			break
		}
	}

	if keeper.Count <= 0 {
		return 0, errors.New("no available id")
	}

	var retId uint64 = 0
	keeper.Count--

	if keeper.FormerEnd == 0 {

		// 位于单一分区, [Start, End)
		// 正常取值，后移
		retId = keeper.Start
		keeper.Start++

	} else if keeper.FormerEnd != 0 && keeper.Start < keeper.FormerEnd {

		// 位于双分区中的[Start, FormerEnd)
		// 正常取值，后移
		retId = keeper.Start
		keeper.Start++

	} else if keeper.FormerEnd != 0 && keeper.Start >= keeper.FormerEnd {

		// 切换到[NewStart, End), Start指针移动到双分区中的[NewStart, End)
		keeper.Start = keeper.NewStart
		glog.Info("switched to id range: ", keeper.NewStart, keeper.End)

		// 正常取值, 后移
		retId = keeper.Start
		keeper.Start++

		// 删除已用尽的[Start, FormerEnd)
		keeper.FormerEnd = 0
		keeper.NewStart = 0
	}

	glog.Info(fmt.Sprintf("id status: %d, [%d,%d), %d", retId,
		keeper.Start, keeper.End, keeper.Count))

	return retId, nil
}

// 后台协程，根据当前可用ID的情况，来获取ID
func (keeper *IdKeeper) refreshId() {
	// 定期执行检查
	routine.ExecRoutineWithInterval("refreshId",
		func() {
			keeper.reallocIdIfNeed()
		},
		time.Second*time.Duration(ReallocCheckInterval))
}

// 扩大可用ID范围
func (keeper *IdKeeper) reallocIdIfNeed() {
	if len(keeper.KeyName) <= 0 {
		glog.Warning("empty key name")
		return
	}

	keeper.Lock.Lock()
	defer keeper.Lock.Unlock()

	// 当可用范围小于阈值时，扩展可用范围
	if keeper.Count > ReallocThreshold {
		return
	}

	// 扩展可用范围
	cmd := redistool.DefaultRedis().IncrBy(context.Background(), keeper.KeyName, ReallocStep)
	val, err := cmd.Result()
	if err != nil {
		glog.Warning("failed to incby ID step: ", err.Error())
		return
	}

	glog.Info("redis incrby id return: ", val)

	// 更新可用范围
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
