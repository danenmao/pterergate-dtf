package serverhelper

import "github.com/danenmao/pterergate-dtf/internal/msgsigner"

func Setup() {
	msgsigner.Setup()
}

func Teardown() {
	msgsigner.Teardown()
}
