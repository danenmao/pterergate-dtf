package idtool

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Intersection_e(t *testing.T) {

	Convey("test two empty id slice", t, func() {

		a := []uint64{}
		b := []uint64{}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should be empty", func() {
			So(c, ShouldBeEmpty)
		})
	})
}
