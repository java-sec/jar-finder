package cmd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	dir = "D:\\java_workspace\\crawler-log-monitor\\out\\artifacts\\crawler_log_monitor_jar"
	pom = "./test_result/pom.xml"
	err := Find(context.Background())
	assert.Nil(t, err)
}
