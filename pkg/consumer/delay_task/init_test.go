package delay_task

import (
	"testing"
	"time"

	"github.com/qml-123/AppService/pkg/redis"
)

func TestAddTask(t *testing.T) {
	_ = redis.InitRedis()
	type args struct {
		taskKey string
		args    []byte
		delay   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				taskKey: "test",
				args: []byte("hello"),
				delay: 1 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddTask(tt.args.taskKey, tt.args.args, tt.args.delay); (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
