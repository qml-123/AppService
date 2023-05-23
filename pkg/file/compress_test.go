package file

import (
	"context"
	"reflect"
	"testing"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/AppService/pkg/log"
	"github.com/qml-123/AppService/pkg/redis"
)

func InitTest() {
	var err error
	if err = log.InitLogger([]string{"http://114.116.15.130:9200"}); err != nil {
		panic(err)
	}

	if err = db.InitDB(); err != nil {
		panic(err)
	}

	if err = redis.InitRedis(); err != nil {
		panic(err)
	}

	if err = id.InitGen(); err != nil {
		panic(err)
	}
}

func TestCompress(t *testing.T) {
	InitTest()
	type args struct {
		ctx       context.Context
		file_keys []string
	}
	tests := []struct {
		name               string
		args               args
		wantFailedFileKeys []string
	}{
		// TODO: Add test cases.
		{
			name: "9026",
			args: args{
				ctx:       context.Background(),
				file_keys: []string{"file_4Sg8J3bj39Y"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFailedFileKeys := Compress(tt.args.ctx, tt.args.file_keys, "./app/"); !reflect.DeepEqual(gotFailedFileKeys, tt.wantFailedFileKeys) {
				t.Errorf("Compress() = %v, want %v", gotFailedFileKeys, tt.wantFailedFileKeys)
			}
		})
	}
}
