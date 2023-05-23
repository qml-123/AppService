package delay_task

import (
	"context"

	"github.com/qml-123/app_log/logger"
)

type testHandle struct{}

func TestHandleMessage(ctx context.Context, msg []byte) error {
	logger.Info(ctx, "get test message: %s", msg)
	return nil
}
