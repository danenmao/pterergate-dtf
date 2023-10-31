package redistool

import (
	"context"
	"errors"
	"time"

	"github.com/golang/glog"
)

const sleepInterval = 50
const defaultExpire = time.Millisecond * 1000 * 20

func Lock(lockName string, timeoutMS uint) error {
	totalCount := timeoutMS / sleepInterval
	var counter uint = 0

	for {
		// 尝试获取锁
		err := tryToLock(lockName, defaultExpire)
		if err == nil {
			return nil
		}

		counter++
		if counter > totalCount {
			return errors.New("timeout to get lock")
		}

		// 等待一段时间再尝试
		time.Sleep(time.Millisecond * sleepInterval)
	}
}

func LockWithExpire(lockName string, timeoutMS uint, expire time.Duration) error {
	totalCount := timeoutMS / sleepInterval
	var counter uint = 0

	for {
		// 尝试获取锁
		err := tryToLock(lockName, expire)
		if err == nil {
			return nil
		}

		counter++
		if counter > totalCount {
			return errors.New("timeout to get lock")
		}

		// 等待一段时间再尝试
		time.Sleep(time.Millisecond * sleepInterval)
	}
}

func Unlock(lockName string) error {
	cmd := DefaultRedis().Del(context.Background(), lockName)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to del lock: ", lockName)
		return err
	}

	glog.Info("succeeded to release lock: ", lockName)

	return nil
}

func RenewLock(lockName string, expire time.Duration) error {
	cmd := DefaultRedis().SetNX(context.Background(), lockName, 1, expire)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to set lock: ", lockName, err)
		return err
	}

	gotLock := cmd.Val()
	if gotLock {
		glog.Warning("lock to be renewed has no owner: ", lockName)
	}

	glog.Info("succeeded to renew lock: ", lockName)
	return nil
}

func tryToLock(lockName string, expire time.Duration) error {
	cmd := DefaultRedis().SetNX(context.Background(), lockName, 1, expire)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to set lock: ", lockName, err)
		return err
	}

	gotLock := cmd.Val()
	if !gotLock {
		glog.Info("lock owned by other: ", lockName)
		return errors.New("lock owned by other")
	}

	glog.Info("got lock: ", lockName)
	return nil
}
