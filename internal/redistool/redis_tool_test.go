package redistool

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_TryToOwnElements_Success(t *testing.T) {
	keyName := "testkey"
	src := []uint64{1, 2, 3}
	owned := []uint64{}

	ClientMock.ExpectTxPipeline()
	for _, v := range src {
		ClientMock.ExpectZRem(keyName, v).SetVal(1)
	}
	ClientMock.ExpectTxPipelineExec()

	err := TryToOwnElements(keyName, &src, &owned)

	Convey("own some elements successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be true", func() {
			So(reflect.DeepEqual(&src, &owned), ShouldBeTrue)
		})
	})
}

func Test_TryToOwnElements_ZRemFail(t *testing.T) {
	keyName := "testkey"
	src := []uint64{1, 2, 3}
	owned := []uint64{}

	ClientMock.ExpectTxPipeline()
	for _, v := range src {
		ClientMock.ExpectZRem(keyName, v).SetVal(0)
	}
	ClientMock.ExpectTxPipelineExec().SetErr(errors.New("exec failed"))

	err := TryToOwnElements(keyName, &src, &owned)

	Convey("failed to own some elements", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
		Convey("should be empty", func() {
			So(len(owned), ShouldBeZeroValue)
		})
	})
}

func Test_GetTimeoutElements_Success(t *testing.T) {
	keyName := "testkey"
	count := uint(3)
	src := []uint64{1, 2, 3}
	result := []uint64{}

	strList := []string{}
	for _, v := range src {
		strList = append(strList, strconv.Itoa(int(v)))
	}

	opt := redis.ZRangeBy{
		Min: "-inf", Max: `^[0-9]{10}$`,
		Offset: 0, Count: int64(count),
	}

	ClientMock.Regexp().ExpectZRangeByScore(keyName, &opt).SetVal(strList)
	err := GetTimeoutElements(keyName, count, &result)

	Convey("get timeout elements successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be true", func() {
			So(reflect.DeepEqual(&src, &result), ShouldBeTrue)
		})
	})
}

func Test_GetTimeoutElements_Fail(t *testing.T) {
	keyName := "testkey"
	count := uint(3)
	result := []uint64{}

	opt := redis.ZRangeBy{
		Min: "-inf", Max: `^[0-9]{10}$`,
		Offset: 0, Count: int64(count),
	}

	ClientMock.Regexp().ExpectZRangeByScore(keyName, &opt).SetErr(errors.New("fail"))
	err := GetTimeoutElements(keyName, count, &result)

	Convey("failed to get timeout elements", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
		Convey("should be true", func() {
			So(len(result) == 0, ShouldBeTrue)
		})
	})
}

func Test_GetTimeoutElements_InvalidInteger(t *testing.T) {
	keyName := "testkey"
	count := uint(3)
	src := []uint64{1, 2, 3}
	result := []uint64{}

	strList := []string{}
	for _, v := range src {
		strList = append(strList, strconv.Itoa(int(v)))
	}

	strList[0] = "invalid"
	expectedList := []uint64{2, 3}

	opt := redis.ZRangeBy{
		Min: "-inf", Max: `^[0-9]{10}$`,
		Offset: 0, Count: int64(count),
	}

	ClientMock.Regexp().ExpectZRangeByScore(keyName, &opt).SetVal(strList)
	ClientMock.ExpectZRem(keyName, "invalid").SetVal(1)

	err := GetTimeoutElements(keyName, count, &result)
	Convey("get timeout elements successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be true", func() {
			So(reflect.DeepEqual(&expectedList, &result), ShouldBeTrue)
		})
	})
}

func Test_GetTimeoutElements_None(t *testing.T) {
	keyName := "testkey"
	count := uint(3)
	result := []uint64{}

	opt := redis.ZRangeBy{
		Min: "-inf", Max: `^[0-9]{10}$`,
		Offset: 0, Count: int64(count),
	}

	ClientMock.Regexp().ExpectZRangeByScore(keyName, &opt).SetVal([]string{})
	err := GetTimeoutElements(keyName, count, &result)

	Convey("get timeout elements successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("should be true", func() {
			So(len(result) == 0, ShouldBeTrue)
		})
	})
}
