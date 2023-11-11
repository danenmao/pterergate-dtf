package redistool

import (
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	fmt.Println("setup...")
	Setup()

	retCode := m.Run()

	fmt.Println("teardown...")
	Teardown()
	os.Exit(retCode)
}

func Test_Lock_Success(t *testing.T) {
	lockName := "test_lock"
	nx := ClientMock.ExpectSetNX(lockName, 1, defaultExpire)
	nx.SetVal(true)

	err := Lock(lockName, 100)

	del := ClientMock.ExpectDel(lockName)
	del.SetVal(1)
	unerr := Unlock(lockName)

	Convey("get a lock successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(unerr, ShouldBeNil)
		})
	})
}

func Test_Lock_LockOwnedByOther(t *testing.T) {
	lockName := "test_lock"
	exp := ClientMock.ExpectSetNX(lockName, 1, defaultExpire)
	exp.SetVal(false)

	err := Lock(lockName, 100)

	Convey("failed to get a lock", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_Unlock_Fail(t *testing.T) {
	lockName := "test_lock"
	exp := ClientMock.ExpectDel(lockName)
	exp.SetErr(errors.New("failed to unlock"))

	err := Unlock(lockName)

	Convey("failed to unlock a lock", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_RenewLock_ValidLock(t *testing.T) {
	lockName := "test_lock"
	exp := ClientMock.ExpectSetNX(lockName, 1, defaultExpire)
	exp.SetVal(false)

	err := RenewLock(lockName, defaultExpire)

	Convey("renew a lock successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
	})
}

func Test_RenewLock_InvalidLock(t *testing.T) {
	lockName := "test_lock"
	exp := ClientMock.ExpectSetNX(lockName, 1, defaultExpire)
	exp.SetVal(true)

	err := RenewLock(lockName, defaultExpire)

	Convey("renew a lock successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
		})
	})
}

func Test_RenewLock_Fail(t *testing.T) {
	lockName := "test_lock"
	exp := ClientMock.ExpectSetNX(lockName, 1, defaultExpire)
	exp.SetErr(errors.New("failed to renew lock"))

	err := RenewLock(lockName, defaultExpire)

	Convey("failed to renew a lock", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
