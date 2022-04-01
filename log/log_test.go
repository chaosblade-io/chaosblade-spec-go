package log

import (
	"context"
	"testing"
)

func Test_info(t *testing.T) {
	Infof(context.WithValue(context.Background(), "uid", "123"), "cpu used: %d", 10)
}
