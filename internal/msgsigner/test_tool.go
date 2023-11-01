package msgsigner

import (
	"os"
	"strings"
)

var gs_wd string

func Setup() {
	gs_wd, _ := os.Getwd()
	idx := strings.LastIndex(gs_wd, "pterergate-dtf")
	rootDir := gs_wd[0 : idx+len("pterergate-dtf")]
	os.Chdir(rootDir)

	KeyPath = "./test/testdata/key.conf"
}

func Teardown() {
	os.Chdir(gs_wd)
}
