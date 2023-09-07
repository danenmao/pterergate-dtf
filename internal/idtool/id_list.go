package idtool

import (
	"sort"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"pterergate-dtf/internal/dbdef"
)

// 将id列表变成分号分隔的数字列表字符串
func IdListToStr(idList *[]uint64) string {
	// 生成客户检测项ID列表
	strList := []string{}
	for _, id := range *idList {
		strList = append(strList, strconv.FormatUint(id, 10))
	}

	// 拼装成字符串
	return strings.Join(strList, dbdef.DBListSeperator)
}

// 将分号分隔的数字列表字符串转换为id列表
func IdStrToIdList(idStr string, idList *[]uint32) error {
	strList := strings.Split(idStr, dbdef.DBListSeperator)
	for _, str := range strList {
		id, err := strconv.Atoi(str)
		if err != nil {
			glog.Warning("failed to convert ", id, err.Error())
			continue
		}

		*idList = append(*idList, uint32(id))
	}

	return nil
}

// 将分号分隔的数字列表字符串转换为id列表
func IdStrToId64List(idStr string, idList *[]uint64) error {
	strList := strings.Split(idStr, dbdef.DBListSeperator)
	for _, str := range strList {
		id, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			glog.Warning("failed to convert ", id, err.Error())
			continue
		}

		*idList = append(*idList, id)
	}

	return nil
}

// 交集
func Intersection(a *[]uint64, b *[]uint64, dest *[]uint64) {

	sort.Slice((*a), func(i, j int) bool { return (*a)[i] < (*a)[j] })
	sort.Slice((*b), func(i, j int) bool { return (*b)[i] < (*b)[j] })

	for _, policy := range *a {

		pos := sort.Search(len(*b), func(i int) bool {
			return (*b)[i] >= policy
		})

		if pos >= len(*b) {
			continue
		}

		if (*b)[pos] != policy {
			continue
		}

		*dest = append(*dest, policy)
	}
}
