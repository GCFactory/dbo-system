package test

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"testing"
)

var (
	testCfgUC = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func Test(t *testing.T) {
	// TODO: реализовать
}
