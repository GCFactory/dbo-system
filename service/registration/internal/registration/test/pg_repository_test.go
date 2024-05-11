package test

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"testing"
)

var (
	testCfg = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func testRepositoryCreateSaga(t *testing.T) {

}
