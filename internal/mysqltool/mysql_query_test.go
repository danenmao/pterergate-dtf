package mysqltool

import (
	"errors"
	"testing"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ReadFromDBByPageCustom_QueryRetryCount(t *testing.T) {
	Convey("test ReadFromDBByPageCustom max query retry limit", t, func() {

		retryCount := 0
		queryFn := func(offset int, limit int) (*sqlx.Rows, error) {
			retryCount += 1
			return nil, errors.New("internal error")
		}

		readFn := func(*sqlx.Rows) error { return nil }

		err := ReadDBByPageWithLimit(queryFn, readFn, 100)

		Convey("return err", func() {
			So(err, ShouldNotBeNil)
		})

		Convey("retry count", func() {
			So(retryCount, ShouldEqual, 3)
		})
	})
}

func Test_ReadFromDBByPageCustom_0Limit(t *testing.T) {
	Convey("test ReadFromDBByPageCustom 0 limit", t, func() {

		queryFn := func(offset int, limit int) (*sqlx.Rows, error) {
			return nil, errors.New("internal error")
		}

		readFn := func(*sqlx.Rows) error { return nil }

		err := ReadDBByPageWithLimit(queryFn, readFn, 0)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})
	})
}

func Test_ReadFromDBByPageCustom_NoData(t *testing.T) {
	Convey("test ReadFromDBByPageCustom no data", t, func() {

		Setup()
		defer Teardown()

		queryFn := func(offset int, limit int) (*sqlx.Rows, error) {
			sql := "select c1 from t"
			mockRows := DBMock.NewRows([]string{"c1"})
			DBMock.ExpectQuery(sql).WillReturnRows(mockRows)
			rows, err := DefaultMySQL().Queryx(sql)
			return rows, err
		}

		readCount := 0
		readFn := func(*sqlx.Rows) error {
			readCount += 1
			return nil
		}

		err := ReadDBByPageWithLimit(queryFn, readFn, 10)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})
	})
}

func Test_ReadFromDBByPageCustom_LessThanDefaultPage(t *testing.T) {

	Convey("test ReadFromDBByPageCustom query data less than one page ", t, func() {
		Setup()
		defer Teardown()

		queryCount := 0
		dataCount := 9
		queryFn := genQueryFn(&dataCount, &queryCount)

		readCount := 0
		readFn := func(*sqlx.Rows) error {
			readCount += 1
			return nil
		}

		err := ReadDBByPageWithLimit(queryFn, readFn, 9)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("query one time", func() {
			So(queryCount, ShouldEqual, 1)
		})

		Convey("read 9 elements", func() {
			So(readCount, ShouldEqual, 9)
		})
	})
}

func Test_ReadFromDBByPageCustom_ReadFailed(t *testing.T) {

	Convey("test ReadFromDBByPageCustom failed to read ", t, func() {
		Setup()
		defer Teardown()

		dataCount := 10
		queryCount := 0
		queryFn := genQueryFn(&dataCount, &queryCount)

		readCount := 0
		readFn := func(rows *sqlx.Rows) error {
			readCount += 1
			return errors.New("internal error")
		}

		err := ReadDBByPageWithLimit(queryFn, readFn, dataCount)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("query one time", func() {
			So(queryCount, ShouldEqual, 1)
		})

		Convey("read 1 elements", func() {
			So(readCount, ShouldEqual, 1)
		})
	})
}

func genQueryFn(dataCount *int, retQueryCount *int) QueryFn {
	return func(offset int, limit int) (*sqlx.Rows, error) {
		*retQueryCount += 1

		thisCount := 0
		if *dataCount > limit {
			thisCount = limit
			*dataCount -= limit
		} else {
			thisCount = *dataCount
			*dataCount = 0
		}

		mockRows := DBMock.NewRows([]string{"c1"})
		for i := 0; i < thisCount; i++ {
			mockRows.AddRow(i)
		}

		sql := "select c1 from t"
		DBMock.ExpectQuery(sql).WillReturnRows(mockRows)
		rows, err := DefaultMySQL().Queryx(sql)
		return rows, err
	}
}

func Test_ReadFromDBByPageCustom_MoreThanOnePage(t *testing.T) {

	Convey("test ReadFromDBByPageCustom query more than one page data ", t, func() {
		Setup()
		defer Teardown()

		totalCount := 101
		dataCount := totalCount
		queryCount := 0
		queryFn := genQueryFn(&dataCount, &queryCount)

		readCount := 0
		readFn := func(rows *sqlx.Rows) error {
			readCount += 1
			return nil
		}

		err := ReadDBByPageWithLimit(queryFn, readFn, totalCount)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("query two time", func() {
			So(queryCount, ShouldEqual, 2)
		})

		Convey("read 101 elements", func() {
			So(readCount, ShouldEqual, totalCount)
		})
	})
}

func Test_ReadFromDBByPageCustom_TwoPage(t *testing.T) {

	Convey("test ReadFromDBByPageCustom query two pages data ", t, func() {
		Setup()
		defer Teardown()

		totalCount := 200
		dataCount := totalCount
		queryCount := 0
		queryFn := genQueryFn(&dataCount, &queryCount)

		readCount := 0
		readFn := func(rows *sqlx.Rows) error {
			readCount += 1
			return nil
		}

		err := ReadDBByPageWithLimit(queryFn, readFn, totalCount)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("query two time", func() {
			So(queryCount, ShouldEqual, 2)
		})

		Convey("read 200 elements", func() {
			So(readCount, ShouldEqual, totalCount)
		})
	})
}

func Test_ReadFromDBByPage_OnePage(t *testing.T) {

	Convey("test ReadFromDBByPage query one pages data ", t, func() {
		Setup()
		defer Teardown()

		totalCount := 1000
		dataCount := totalCount
		queryCount := 0
		queryFn := genQueryFn(&dataCount, &queryCount)

		readCount := 0
		readFn := func(rows *sqlx.Rows) error {
			readCount += 1
			return nil
		}

		err := ReadDBByPage(queryFn, readFn)

		Convey("return nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("query 10 time", func() {
			So(queryCount, ShouldEqual, 10)
		})

		Convey("read 1000 elements", func() {
			So(readCount, ShouldEqual, totalCount)
		})
	})
}

func Test_AssembleListSQL_IntList(t *testing.T) {
	Convey("test AssembleListSQL int list", t, func() {

		part1Sql := "select a from t where b in ("
		part2Sql := ") limit 100"
		part1Args := []interface{}{}
		listArgs := []uint64{1, 2, 3}
		part2Args := []interface{}{}

		sql := ""
		resultArgs := []interface{}{}
		AssembleInRangeSQL(part1Sql, part2Sql, &part1Args, &listArgs, &part2Args, &sql, &resultArgs)

		Convey("sql shoud be equal", func() {
			So(sql, ShouldEqual, "select a from t where b in (?,?,?) limit 100")
		})

		Convey("args shoud be equal", func() {
			So(len(resultArgs), ShouldEqual, len(listArgs))
		})
	})
}

func Test_AssembleListSQLTemplate_StrList(t *testing.T) {
	Convey("test AssembleListSQLTemplate string list", t, func() {

		part1Sql := "select a from t where b in ("
		part2Sql := ") limit 100"
		part1Args := []interface{}{}
		listArgs := []interface{}{"1", "2", "3"}
		part2Args := []interface{}{}

		sql := ""
		resultArgs := []interface{}{}
		AssembleInRangeSQLTemplate(part1Sql, part2Sql, &part1Args, &listArgs, &part2Args, &sql, &resultArgs)

		Convey("sql shoud be equal", func() {
			So(sql, ShouldEqual, "select a from t where b in (?,?,?) limit 100")
		})

		Convey("args shoud be equal", func() {
			So(len(resultArgs), ShouldEqual, len(listArgs))
		})
	})
}

func Test_AssemberInsertOrUpdateSQL_3Elem(t *testing.T) {
	Convey("test AssemberInsertOrUpdateSQL ", t, func() {

		part1Sql := "insert into t (a,b,c) values "
		fieldSql := " (:a, :b, :c) "
		part2Sql := " (?,?,?) "
		part3Sql := " on duplicate key update a=values(a), b=values(b), c=values(c)"

		sql := ""
		AssemberInsertOrUpdateSQL(part1Sql, part2Sql, part3Sql, fieldSql, 3, &sql)

		Convey("sql should be eqaul", func() {
			So(sql, ShouldEqual, "insert into t (a,b,c) values  (:a, :b, :c) , (?,?,?) , (?,?,?) "+
				" on duplicate key update a=values(a), b=values(b), c=values(c)")
		})

	})
}

func Test_AssemberInsertOrUpdateSQL_0Elem(t *testing.T) {
	Convey("test AssemberInsertOrUpdateSQL with 0 elem ", t, func() {

		part1Sql := "insert into t (a,b,c) values "
		fieldSql := " (:a, :b, :c) "
		part2Sql := " (?,?,?) "
		part3Sql := " on duplicate key update a=values(a), b=values(b), c=values(c)"

		sql := ""
		err := AssemberInsertOrUpdateSQL(part1Sql, part2Sql, part3Sql, fieldSql, 0, &sql)

		Convey("error should occur", func() {
			So(err, ShouldNotBeNil)
		})

	})
}
