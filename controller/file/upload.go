package file

import (
	"bytes"
	"context"
	"encoding/binary"
	"os"
)

func Upload(ctx context.Context, user_id string, file []byte) (string, error) {
	var err error
	//ori_file, err := ffmpeg.DecodeAV1(file)
	if err != nil {
		return "", err
	}
	fp, err := os.Create("new_file")
	if err != nil {
		return "", nil
	}
	defer fp.Close()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, file)
	fp.Write(buf.Bytes())
	return "new_file", nil
}
