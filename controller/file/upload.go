package file

import (
	"bytes"
	"context"
	"encoding/binary"
	"os"

	"github.com/qml-123/AppService/cgo/av1"
)

func Upload(ctx context.Context, user_id string, file []byte) (string, error) {
	var err error
	ori_file, err := av1.AV1Decode(file)
	if err != nil {
		return "", err
	}
	fp, err := os.Create("new_file")
	if err != nil {
		return "", nil
	}
	defer fp.Close()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ori_file)
	fp.Write(buf.Bytes())
	return "new_file", nil
}