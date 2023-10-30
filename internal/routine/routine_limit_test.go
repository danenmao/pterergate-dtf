package routine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Count_0(t *testing.T) {
	limiter := CountLimiter{UpperLimit: 10}

	Convey("test initial count is zero", t, func() {
		Convey("counter should be zero", func() {
			So(limiter.Count(), ShouldBeZeroValue)
		})
	})
}

func Test_Count_IncrToOne(t *testing.T) {
	limiter := CountLimiter{UpperLimit: 10}
	limiter.Incr()

	Convey("test the increased count is one", t, func() {
		Convey("counter should be one", func() {
			So(limiter.Count(), ShouldEqual, 1)
			So(limiter.IsFull(), ShouldBeFalse)
		})
	})
}

func Test_Count_IncrManyTimes(t *testing.T) {
	limiter := CountLimiter{UpperLimit: 10}
	for i := 0; i < 9; i++ {
		limiter.Incr()
	}

	Convey("test the increased count is nine", t, func() {
		Convey("counter should be nine", func() {
			So(limiter.Count(), ShouldEqual, 9)
			So(limiter.IsFull(), ShouldBeFalse)
		})
	})
}

func Test_IncrIfNotFull_IncrToUpperLimit(t *testing.T) {
	const UPPERLIMIT = 10
	limiter := CountLimiter{UpperLimit: UPPERLIMIT}
	for i := 0; i < UPPERLIMIT+1; i++ {
		limiter.IncrIfNotFull()
	}

	Convey("test the increased count is the upper limit", t, func() {
		Convey("counter should be the upper limit", func() {
			So(limiter.Count(), ShouldEqual, UPPERLIMIT)
			So(limiter.IsFull(), ShouldBeTrue)
		})
	})
}

func Test_Decr_DecrManyTimes(t *testing.T) {
	const UPPERLIMIT = 10
	limiter := CountLimiter{UpperLimit: UPPERLIMIT}
	for i := 0; i < UPPERLIMIT; i++ {
		limiter.IncrIfNotFull()
		limiter.Decr()
	}

	Convey("test the decreased count is zero", t, func() {
		Convey("counter should be zero", func() {
			So(limiter.Count(), ShouldBeZeroValue)
		})
	})
}
