package redistool

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
)

// 试图获取元素的所有权
func TryToOwnElemList(keyName string, elemList *[]uint64, ownedElemList *[]uint64) error {

	// 从 key 中删除子任务
	pipeline := DefaultRedis().Pipeline()
	for _, elem := range *elemList {
		pipeline.ZRem(context.Background(), keyName, elem)
	}

	// 执行pipeline
	cmdList, err := pipeline.Exec(context.Background())
	if err != nil {
		glog.Warning("failed to exec pipeline: ", err)
		return err
	}

	// 检查删除结果，删除成功则拥有该任务的处理所有权
	for idx, cmd := range cmdList {

		// 取结果对应的任务的ID
		elem := (*elemList)[idx]

		// 检查ZRem id的结果
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
			//glog.Info("elem owned by other: ", elem)
			continue
		}

		// 为拥有所有权的任务
		*ownedElemList = append(*ownedElemList, elem)
		glog.Info("owned completed elem: ", elem)
	}

	return nil
}

// 获取超时的元素列表
func GetTimeoutElemList(keyName string, count uint, elemList *[]uint64) error {

	if elemList == nil {
		panic("invalid elem list pointer")
	}

	// 从zset中取超时的子任务
	now := time.Now().Unix()
	nowStr := strconv.FormatUint(uint64(now), 10)
	opt := redis.ZRangeBy{
		Min: "-inf", Max: nowStr,
		Offset: 0, Count: int64(count),
	}

	cmd := DefaultRedis().ZRangeByScore(context.Background(), keyName, &opt)
	err := cmd.Err()
	if err != nil {
		glog.Warning("failed to get timeout elem from redis: ", err)
		return err
	}

	strList := cmd.Val()
	if len(strList) > 0 {
		glog.Info("got timeout elem: ", strList)
	}

	// 转换查询到的子任务ID
	var wrongFormatList = []interface{}{}
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			// 如果转换失败，说明数据格式错误，移除元素
			glog.Warning("failed to convert timeout elem id: ", str)
			wrongFormatList = append(wrongFormatList, str)
			continue
		}

		*elemList = append(*elemList, id)
	}

	// 删除转换失败的任务数据
	DefaultRedis().ZRem(context.Background(), keyName, wrongFormatList...)

	// 如果列表为空，表示没有超时的任务
	if len(*elemList) == 0 {
		//glog.Info("get empty list, no timeout elem")
		return nil
	}

	glog.Info("got timeout elems: ", *elemList)
	return nil
}
