package webrtc

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type someStruct struct {
	value int
}

func TestCGOMap(t *testing.T) {
	Convey("CGOMap", t, func() {
		Convey("Gets and Sets correctly", func() {
			m := NewCGOMap()

			i1 := m.Set(someStruct{123})
			So(i1, ShouldEqual, 1)
			So(m.Get(1).(someStruct).value, ShouldEqual, 123)

			i2 := m.Set(someStruct{456})
			So(i2, ShouldEqual, 2)
			So(m.Get(2).(someStruct).value, ShouldEqual, 456)
		})

		Convey("Deletes correctly", func() {
			m := NewCGOMap()
			i := m.Set(someStruct{234})
			So(m.Get(i).(someStruct).value, ShouldEqual, 234)
			m.Delete(i)
			So(m.pointers[i], ShouldBeNil)
		})
	})
}
