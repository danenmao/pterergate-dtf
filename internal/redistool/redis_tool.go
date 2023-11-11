package redistool

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
)

// try to get the ownership of elements
func TryToOwnElements(keyName string, srcElements *[]uint64, ownedElements *[]uint64) error {
	// remove elements
	pipeline := DefaultRedis().TxPipeline()
	for _, elem := range *srcElements {
		pipeline.ZRem(context.Background(), keyName, elem)
	}

	// execute the pipeline
	cmdList, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	// check if get the ownership of elements
	for idx, cmd := range cmdList {
		elem := (*srcElements)[idx]

		// check the result of ZRem
		intCmd, ok := cmd.(*redis.IntCmd)
		if !ok {
			glog.Warning("failed to convert cmd: ", elem, cmd)
			continue
		}

		err = intCmd.Err()
		if err != nil {
			glog.Info("failed to zrem elem from list:", elem, ",", err)
			continue
		}

		result := intCmd.Val()
		if result == 0 {
			continue
		}

		// get the ownership of an element
		*ownedElements = append(*ownedElements, elem)
		glog.Info("owned an element: ", elem)
	}

	return nil
}

func GetTimeoutElements(keyName string, count uint, retElements *[]uint64) error {
	if retElements == nil {
		panic("invalid element list pointer")
	}

	// get timeout elements
	now := time.Now().Unix()
	nowStr := strconv.FormatUint(uint64(now), 10)
	opt := redis.ZRangeBy{
		Min: "-inf", Max: nowStr,
		Offset: 0, Count: int64(count),
	}

	cmd := DefaultRedis().ZRangeByScore(context.Background(), keyName, &opt)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get timeout element from redis: ", err)
		return err
	}

	strList := cmd.Val()
	if len(strList) > 0 {
		glog.Info("got a timeout element: ", strList)
	}

	// covert timeout elements
	var wrongFormatList = []interface{}{}
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			// invalid digital format，remove it
			glog.Warning("failed to convert a timeout element id: ", str)
			wrongFormatList = append(wrongFormatList, str)
			continue
		}

		*retElements = append(*retElements, id)
	}

	// remove the elements that failed to covert
	if len(wrongFormatList) > 0 {
		DefaultRedis().ZRem(context.Background(), keyName, wrongFormatList...)
	}

	// no timeout elements
	if len(*retElements) == 0 {
		return nil
	}

	glog.Info("got timeout elements: ", *retElements)
	return nil
}
