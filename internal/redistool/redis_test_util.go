package redistool

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var s_actualClient *goredis.Client
var s_mockClient *goredis.Client
var ClientMock redismock.ClientMock

func Setup() {

	// save actual client
	if gs_RedisClient != nil {
		s_actualClient = gs_RedisClient
	}

	// generate a mock client
	s_mockClient, ClientMock = redismock.NewClientMock()

	// overwrite the default client
	gs_RedisClient = s_mockClient
}

func Teardown() {

	// restore the default client
	if s_actualClient != nil {
		gs_RedisClient = s_actualClient
		s_actualClient = nil
	}

	// release the mock client
	if s_mockClient != nil {
		s_mockClient.Close()
		s_mockClient = nil
	}
}
