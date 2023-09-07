package idtool

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Intersection_BothEmptySlice(t *testing.T) {

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

func Test_Intersection_OneEmptySlice(t *testing.T) {

	Convey("test one empty id slice", t, func() {

		a := []uint64{1, 2, 3, 4}
		b := []uint64{}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should be empty", func() {
			So(c, ShouldBeEmpty)
		})
	})
}

func Test_Intersection_NoSameElementSlice(t *testing.T) {

	Convey("test no same id slice", t, func() {

		a := []uint64{1, 2, 3, 4}
		b := []uint64{5, 6, 7, 8}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should be empty", func() {
			So(c, ShouldBeEmpty)
		})
	})
}

func Test_Intersection_OneSameElementSlice(t *testing.T) {

	Convey("test one empty id slice", t, func() {

		a := []uint64{1, 2, 3, 4}
		b := []uint64{1, 5, 6, 7}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should have one element", func() {
			So(len(c), ShouldEqual, 1)
		})

		Convey("c[0] should be 1", func() {
			So(c[0], ShouldEqual, 1)
		})
	})
}

func Test_Intersection_SameSlice(t *testing.T) {

	Convey("test the same id slice", t, func() {

		a := []uint64{1, 2, 3, 4}
		b := []uint64{1, 2, 3, 4}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should have the same count of element as a", func() {
			So(len(c), ShouldEqual, len(a))
		})

		Convey("c should have the same count of element as b", func() {
			So(len(c), ShouldEqual, len(b))
		})

		Convey("c should be the same as a", func() {
			So(reflect.DeepEqual(a, c), ShouldBeTrue)
		})

		Convey("c should be the same as b", func() {
			So(reflect.DeepEqual(b, c), ShouldBeTrue)
		})
	})
}

func Test_Intersection_TwoSameElementSlice(t *testing.T) {

	Convey("test one empty id slice", t, func() {

		a := []uint64{1, 2, 3, 4}
		b := []uint64{2, 3, 5, 5}
		c := []uint64{}
		Intersection(&a, &b, &c)

		Convey("c should have 2 element", func() {
			So(len(c), ShouldEqual, 2)
		})

		Convey("c should be the same", func() {
			So(reflect.DeepEqual(c, []uint64{2, 3}), ShouldBeTrue)
		})
	})
}

func Test_IdStrToId64List_EmptySlice(t *testing.T) {

	Convey("test empty id slice", t, func() {

		idStr := ""
		idList := []uint64{}
		IdStrToId64List(idStr, &idList)

		Convey("c should have no element", func() {
			So(idList, ShouldBeEmpty)
		})

	})
}

func Test_IdStrToId64List_OnlyHaveSeperators(t *testing.T) {

	Convey("test empty id slice", t, func() {

		idStr := ";;;;"
		idList := []uint64{}
		IdStrToId64List(idStr, &idList)

		Convey("c should have no element", func() {
			So(idList, ShouldBeEmpty)
		})

	})
}

func Test_IdStrToId64List_OnlyOneElement(t *testing.T) {

	Convey("test one id slice", t, func() {

		idStr := "1"
		idList := []uint64{}
		IdStrToId64List(idStr, &idList)

		Convey("c should have 1 element", func() {
			So(reflect.DeepEqual(idList, []uint64{1}), ShouldBeTrue)
		})

	})
}

func Test_IdStrToId64List_OnlyOneElementEndWithSeperator(t *testing.T) {

	Convey("test one id slice ends with seperator", t, func() {

		idStr := "1;"
		idList := []uint64{}
		IdStrToId64List(idStr, &idList)

		Convey("c should have 1 element", func() {
			So(reflect.DeepEqual(idList, []uint64{1}), ShouldBeTrue)
		})

	})
}

func Test_IdStrToId64List_ManyElement(t *testing.T) {

	Convey("test many id slice", t, func() {

		idStr := "1;2;3;4"
		idList := []uint64{}
		IdStrToId64List(idStr, &idList)

		Convey("c should have 1 element", func() {
			So(reflect.DeepEqual(idList, []uint64{1, 2, 3, 4}), ShouldBeTrue)
		})

	})
}

func Test_IdStrToIdList_OnlyOneElement(t *testing.T) {

	Convey("test one id slice", t, func() {

		idStr := "1"
		idList := []uint32{}
		IdStrToIdList(idStr, &idList)

		Convey("c should have 1 element", func() {
			So(reflect.DeepEqual(idList, []uint32{1}), ShouldBeTrue)
		})

	})
}

func Test_IdStrToIdList_EmptySlice(t *testing.T) {

	Convey("test empty id slice", t, func() {

		idStr := ""
		idList := []uint32{}
		IdStrToIdList(idStr, &idList)

		Convey("c should have no element", func() {
			So(idList, ShouldBeEmpty)
		})

	})
}

func Test_IdStrToIdList_InvalidElement(t *testing.T) {

	Convey("test empty id slice", t, func() {

		idStr := "a;b;c;d"
		idList := []uint32{}
		IdStrToIdList(idStr, &idList)

		Convey("c should have no element", func() {
			So(idList, ShouldBeEmpty)
		})

	})
}

func Test_IdListToStr_EmptySlice(t *testing.T) {

	Convey("test empty id slice", t, func() {

		idList := []uint64{}
		idStr := IdListToStr(&idList)

		Convey("str should be empty", func() {
			So(idStr, ShouldBeEmpty)
		})

	})
}

func Test_IdListToStr_OneElement(t *testing.T) {

	Convey("test one id slice", t, func() {

		idList := []uint64{1}
		idStr := IdListToStr(&idList)

		Convey("str should be equal", func() {
			So(idStr, ShouldEqual, "1")
		})

	})
}

func Test_IdListToStr_ManyElement(t *testing.T) {

	Convey("test many id slice", t, func() {

		idList := []uint64{1, 2, 3, 4, 10000, 10000000}
		idStr := IdListToStr(&idList)

		Convey("str should be equal", func() {
			So(idStr, ShouldEqual, "1;2;3;4;10000;10000000")
		})

	})
}
