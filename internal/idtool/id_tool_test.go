package idtool

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Intersection_empty(t *testing.T) {

	Convey("test two empty id slice", t, func() {

		Convey("c should be empty", func() {
			So(nil, ShouldBeEmpty)
		})
	})
}
