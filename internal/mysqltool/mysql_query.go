package mysqltool

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

// 分页读取的函数框架
type QueryFn func(offset int, limit int) (*sqlx.Rows, error)
type ReadRowFn func(*sqlx.Rows) error

// 分页读取MySQL表, 有总记录数限制
func ReadFromDBByPage(queryFn QueryFn, readFn ReadRowFn) error {
	const TotalCountLimit = 1000
	return ReadFromDBByPageCustom(queryFn, readFn, TotalCountLimit)
}

// 分页读取MySQL表，可设置记录数上限
func ReadFromDBByPageCustom(queryFn QueryFn, readFn ReadRowFn, totalCountLimit int) error {
	const CountLimitPerQuery = 100
	const MaxRetryCount = 3
	var TotalCountLimit = totalCountLimit
	var QueryTimesLimit = TotalCountLimit / CountLimitPerQuery

	// 分页查询策略项记录
	totalCount := 0
	for i := 0; i < QueryTimesLimit; i++ {
		startIndex := i * CountLimitPerQuery

		// 查询一页记录，带有失败重试
		var j int
		var needToBreak bool = false
		var querySuccessFlag bool = false
		for j = 0; j < MaxRetryCount; j++ {
			rows, err := queryFn(startIndex, CountLimitPerQuery)
			//defer rows.Close()

			// 查询失败会进行重试
			if err != nil {
				glog.Warning("failed to query item: ", err.Error())
				continue
			}

			itemCountInPage := 0
			for rows.Next() {
				err = readFn(rows)
				if err != nil {
					glog.Warning("failed to scan item: ", err.Error())
					break
				}

				totalCount += 1
				itemCountInPage += 1
			}

			if err = rows.Err(); err != nil {
				glog.Warning("query record error: ", err)
			}

			querySuccessFlag = true

			// 当前页面内的数量小于上限，查询完成，跳出
			if itemCountInPage < CountLimitPerQuery {
				needToBreak = true
				break
			}

			// 本批次读取、解析成功，无需重试
			break
		}

		// 超过重试上限且最后一次查询未成功
		if j >= MaxRetryCount && !querySuccessFlag {
			return errors.New("failure time exceeds max retry times")
		}

		if needToBreak {
			break
		}
	}

	return nil
}

// 组装 in () 类型的SQL语句及参数
func AssembleListSQL(
	part1Sql string, part2Sql string,
	part1Args *[]interface{}, listArgs *[]uint64, part2Args *[]interface{},
	sqlRet *string, fieldArgs *[]interface{},
) {

	*fieldArgs = append(*fieldArgs, *part1Args...)

	var sqlPart []string
	for _, id := range *listArgs {
		*fieldArgs = append(*fieldArgs, id)
		sqlPart = append(sqlPart, "?")
	}

	*fieldArgs = append(*fieldArgs, *part2Args...)

	fieldSql := strings.Join(sqlPart, ",")
	sql := fmt.Sprintf("%s%s%s",
		part1Sql,
		fieldSql,
		part2Sql,
	)

	*sqlRet = sql
}

// 组装 in () 类型的SQL语句及参数
func AssembleListSQLTemplate(
	part1Sql string, part2Sql string,
	part1Args *[]interface{}, listArgs *[]interface{}, part2Args *[]interface{},
	sqlRet *string, fieldArgs *[]interface{},
) {

	*fieldArgs = append(*fieldArgs, *part1Args...)

	var sqlPart []string
	for _, id := range *listArgs {
		*fieldArgs = append(*fieldArgs, id)
		sqlPart = append(sqlPart, "?")
	}

	*fieldArgs = append(*fieldArgs, *part2Args...)

	fieldSql := strings.Join(sqlPart, ",")
	sql := fmt.Sprintf("%s%s%s",
		part1Sql,
		fieldSql,
		part2Sql,
	)

	*sqlRet = sql
}

// 组装 insert ... update 语句
func AssemberInsertOrUpdateSQL(
	part1Sql string, part2Sql string, part3Sql string, fieldSql string,
	elemCount uint32,
	sqlRet *string,
) error {

	if elemCount == 0 {
		glog.Warning("empty item list to update")
		return errors.New("empty list")
	}

	valuesSqlList := []string{}
	valuesSqlList = append(valuesSqlList, fieldSql)
	for i := 1; i < int(elemCount); i++ {
		valuesSqlList = append(valuesSqlList, part2Sql)
	}

	sql := fmt.Sprintf("%s%s%s",
		part1Sql,
		strings.Join(valuesSqlList, ","),
		part3Sql,
	)

	*sqlRet = sql
	return nil
}
