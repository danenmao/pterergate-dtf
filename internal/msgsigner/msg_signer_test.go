package msgsigner

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MsgInit_Success(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	signer := NewMsgSigner()

	Convey("new a MsgSigner successfully", t, func() {
		Convey("should not be nil", func() {
			So(signer, ShouldNotBeNil)
			So(signer.privateKey, ShouldNotBeNil)
			So(signer.publicKey, ShouldNotBeNil)
		})
	})
}

func Test_MsgInit_Fail(t *testing.T) {
	KeyPath = "./test/testdata/invalid.conf"
	signer := NewMsgSigner()

	Convey("failed to new a MsgSigner", t, func() {
		Convey("should not be nil", func() {
			So(signer, ShouldNotBeNil)
			So(signer.privateKey, ShouldBeNil)
			So(signer.publicKey, ShouldBeNil)
		})
	})
}

func Test_Sign_Success(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	signer := NewMsgSigner()
	secret, err := signer.Sign("tester", "test message", []string{"tester"},
		"test message string", time.Minute)

	Convey("sign a message successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(secret, ShouldNotBeNil)
		})
	})
}

func Test_Verify_Success(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	signer := NewMsgSigner()
	msg := "test message string"
	secret, err := signer.Sign("tester", "test message", []string{"tester"},
		msg, time.Minute)
	plain, err := signer.Verify(secret)

	Convey("verify a secret successfully", t, func() {
		Convey("should be nil", func() {
			So(err, ShouldBeNil)
			So(plain, ShouldEqual, msg)
		})
	})
}

func Test_Verify_Expired(t *testing.T) {
	KeyPath = "./test/testdata/key.conf"
	signer := NewMsgSigner()
	msg := "test message string"
	secret, err := signer.Sign("tester", "test message", []string{"tester"},
		msg, 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	plain, err := signer.Verify(secret)

	Convey("verify a expired secret", t, func() {
		Convey("should expired", func() {
			So(err, ShouldNotBeNil)
			So(plain, ShouldEqual, "")
		})
	})
}
