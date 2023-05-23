package av1

import "testing"

func TestConvertToAV1(t *testing.T) {
	type args struct {
		apiKey     string
		inputPath  string
		outputPath string
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
				apiKey: api_key,
				inputPath: "/opt/app/file_4SiubNpow8j/file_4SiubNpow8j.m3u8",
				outputPath: "/opt/app/file_4SiubNpow8j/file_4SiubNpow8j.av1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertToAV1(tt.args.apiKey, tt.args.inputPath, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("ConvertToAV1() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
