package av1

import (
	"context"
	"testing"
)

func TestConvertTsToAv1(t *testing.T) {
	type args struct {
		ctx      context.Context
		m3u8Path string
		tsDir    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				ctx:      context.Background(),
				m3u8Path: "/Users/bytedance/Documents/CLionProject/cpp/app_client/build/file_4SihNifBcZj/file_4SihNifBcZj.m3u8",
				tsDir:    "/Users/bytedance/Documents/CLionProject/cpp/app_client/build/file_4SihNifBcZj",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertTsToAv1(tt.args.ctx, tt.args.m3u8Path, tt.args.tsDir); (err != nil) != tt.wantErr {
				t.Errorf("ConvertTsToAv1() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
