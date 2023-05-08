package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func GetFileData(filePath string) []byte {
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()

	buff := make([]byte, 55) // 55=该文本的长度

	for {
		lens, err := fp.Read(buff)
		if err == io.EOF || lens < 0 {
			break
		}
	}
	return buff
}

func SaveFileData(data []byte, filePath string) error {
	fp, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	fp.Write(buf.Bytes())
	return nil
}
