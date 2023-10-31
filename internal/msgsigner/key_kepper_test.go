package msgsigner

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	fmt.Println("setup...")
	wd, _ := os.Getwd()
	idx := strings.Index(wd, "pterergate-dtf")
	rootDir := wd[0 : idx+len("pterergate-dtf")]
	os.Chdir(rootDir)

	retCode := m.Run()

	fmt.Println("teardown...")
	os.Chdir(wd)
	os.Exit(retCode)
}

func Test_Load_Initial(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	keeper := &KeyKeeper{}
	err := keeper.Load()

	Convey("load key from a file", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldBeNil)
			So(keeper.privateKey, ShouldNotBeNil)
			So(keeper.publicKey, ShouldNotBeNil)
		})
	})
}

func Test_Load_InvalidKeyFormat(t *testing.T) {
	KeyPath = "./test/testdata/invalid.conf"
	keeper := &KeyKeeper{}
	err := keeper.Load()
	k, kerr := keeper.GetPrivateKey()
	pk, pkerr := keeper.GetPublicKey()

	Convey("load key from an invalid format", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldNotBeNil)
			So(keeper.privateKey, ShouldBeNil)
			So(keeper.publicKey, ShouldBeNil)
			So(k, ShouldBeNil)
			So(kerr, ShouldNotBeNil)
			So(pk, ShouldBeNil)
			So(pkerr, ShouldNotBeNil)
		})
	})
}

func Test_Load_InvalidPublicKey(t *testing.T) {
	KeyPath = "./test/testdata/invalidpk.conf"
	keeper := &KeyKeeper{}
	keeper.Load()
	k, err := keeper.GetPrivateKey()
	pk, pkerr := keeper.GetPublicKey()

	Convey("load key from a file", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldBeNil)
			So(pkerr, ShouldNotBeNil)
			So(k, ShouldNotBeNil)
			So(pk, ShouldBeNil)
		})
	})
}

func Test_GetPrivateKey_Normal(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	keeper := &KeyKeeper{}
	keeper.Load()
	k, err := keeper.GetPrivateKey()
	pk, err := keeper.GetPublicKey()

	Convey("load key from a file", t, func() {
		Convey("should not be nil", func() {
			So(err, ShouldBeNil)
			So(k, ShouldNotBeNil)
			So(pk, ShouldNotBeNil)
		})
	})
}
