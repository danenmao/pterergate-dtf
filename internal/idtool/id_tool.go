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
	RenewIDStep        = 100
	RenewThreshold     = 20
	RenewCheckInterval = 5
)

// ID锁
var s_IDLock sync.RWMutex

// ID数据
var (
	s_KeyName string
)

// ID可用范围：[Start, FormerEnd),[NewStart, End)
type IDKeeper struct {
	KeyName   string
	Step      uint32
	Count     uint32
	Start     uint64
	FormerEnd uint64
	NewStart  uint64
	End       uint64
}

var s_AvailableID = IDKeeper{}

// 初始化
func init() {

}

// 初始化
func Init(keyName string) error {

	if len(keyName) <= 0 {
		return errors.New("empty key name")
	}

	s_KeyName = keyName
	s_AvailableID.KeyName = keyName
	s_AvailableID.Step = RenewIDStep

	// TODO: 定期同步至MySQL, 处理迁移的场景
	// 启动后台维护协程
	go getIDRoutine()

	return nil
}

// 获取可用的ID
func GetAvailableId(keyName string) (uint64, error) {

	if keyName != s_KeyName {
		return 0, errors.New("mismatched key name")
	}

	s_IDLock.Lock()
	defer s_IDLock.Unlock()

	// 检查是否有可用ID
	for i := 0; i < 2; i++ {

		if s_AvailableID.Count <= 0 {
			s_IDLock.Unlock()

			glog.Info("renew id range immediately")
			renewAvailableIDIfNeed()

			s_IDLock.Lock()
			continue
		} else {
			break
		}
	}

	if s_AvailableID.Count <= 0 {
		return 0, errors.New("no available id")
	}

	var retId uint64 = 0
	s_AvailableID.Count--

	if s_AvailableID.FormerEnd == 0 {

		// 位于单一分区, [Start, End)
		// 正常取值，后移
		retId = s_AvailableID.Start
		s_AvailableID.Start++

	} else if s_AvailableID.FormerEnd != 0 && s_AvailableID.Start < s_AvailableID.FormerEnd {

		// 位于双分区中的[Start, FormerEnd)
		// 正常取值，后移
		retId = s_AvailableID.Start
		s_AvailableID.Start++

	} else if s_AvailableID.FormerEnd != 0 && s_AvailableID.Start >= s_AvailableID.FormerEnd {

		// 切换到[NewStart, End), Start指针移动到双分区中的[NewStart, End)
		s_AvailableID.Start = s_AvailableID.NewStart
		glog.Info("switched to id range: ", s_AvailableID.NewStart, s_AvailableID.End)

		// 正常取值, 后移
		retId = s_AvailableID.Start
		s_AvailableID.Start++

		// 删除已用尽的[Start, FormerEnd)
		s_AvailableID.FormerEnd = 0
		s_AvailableID.NewStart = 0
	}

	glog.Info(fmt.Sprintf("id status: -> %d, [%d,%d), %d", retId,
		s_AvailableID.Start, s_AvailableID.End, s_AvailableID.Count))

	return retId, nil
}

// 后台协程，根据当前可用ID的情况，来获取ID
func getIDRoutine() {
	// 定期执行检查
	routine.ExecRoutineByDuration("getIDRoutine", renewAvailableIDIfNeed,
		time.Second*time.Duration(RenewCheckInterval))
}

// 扩大可用ID范围
func renewAvailableIDIfNeed() {

	if len(s_KeyName) <= 0 {
		glog.Warning("empty key name")
		return
	}

	s_IDLock.Lock()
	defer s_IDLock.Unlock()

	// 当可用范围小于阈值时，扩展可用范围
	if s_AvailableID.Count > RenewThreshold {
		return
	}

	// 扩展可用范围
	cmd := redistool.DefaultRedis().IncrBy(context.Background(), s_KeyName, RenewIDStep)
	val, err := cmd.Result()
	if err != nil {
		glog.Warning("failed to incby ID step: ", err.Error())
		return
	}

	glog.Info("redis incrby id return: ", val)

	// 更新可用范围
	s_AvailableID.Count += uint32(s_AvailableID.Step)

	// [Start, FormerEnd), [NewStart, End)
	s_AvailableID.FormerEnd = s_AvailableID.End
	s_AvailableID.End = uint64(val) + 1

	if s_AvailableID.Start == 0 {
		s_AvailableID.Start = s_AvailableID.End - uint64(s_AvailableID.Step)
	}

	s_AvailableID.NewStart = s_AvailableID.End - uint64(s_AvailableID.Step)

	if s_AvailableID.FormerEnd > 0 && s_AvailableID.FormerEnd-s_AvailableID.Start > uint64(s_AvailableID.Step) {
		glog.Error("error id range: [", s_AvailableID.Start, s_AvailableID.FormerEnd, ")")
	}

	glog.Info(fmt.Sprintf("renew id range: [%d,%d), [%d,%d), %d",
		s_AvailableID.Start, s_AvailableID.FormerEnd,
		s_AvailableID.NewStart, s_AvailableID.End, s_AvailableID.Count))
}
