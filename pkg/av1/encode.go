package av1

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/qml-123/app_log/logger"
)

type DirFile struct {
	Dir     string
	File    string
	FileKey string
	UserID  int64
}

// 合并
func EncodeSplitFile(ctx context.Context, d *DirFile) error {
	logger.Info(ctx, "EncodeSplitFile begin, file_key: %v, file_name: %v", d.FileKey, d.File)
	var stderr bytes.Buffer
	// 创建 FFmpeg 命令
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		d.File,
		"-vf",
		"vflip", // 添加vflip参数来翻转视频画面
		"-c",
		"copy",
		d.Dir+d.FileKey+".mp4")

	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("stderr: %v, err: %v", stderr.String(), err)
	}
	return nil
}

// Encode each segment independently
func EncodeEachFile(ctx context.Context, d *DirFile) error {
	logger.Info(ctx, "EncodeEachFile begin, file_key: %v, file_name: %v", d.FileKey, d.File)
	movFile := d.Dir + d.FileKey + ".mp4"
	outputFile := d.Dir + d.FileKey + "_output.av1"
	cmd := exec.Command("ffmpeg", "-i", movFile, "-c:v", "libaom-av1", "-strict", "-2", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Concatenate the encoded segments
func Concatenate(ctx context.Context, d *DirFile) error {
	logger.Info(ctx, "Concatenate begin, file_key: %v, file_name: %v", d.FileKey, d.File)
	file, err := os.Create(d.Dir + "list.txt")
	if err != nil {
		return err
	}

	defer file.Close()

	files, err := os.ReadDir(d.Dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), "_encoded.m3u8") {
			continue
		}
		_, err = file.WriteString("file '" + f.Name() + "'\n")
		if err != nil {
			return err
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", d.Dir+"list.txt", "-c", "copy", d.Dir+"final_output.m3u8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("stderr: %v, err: %v", stderr.String(), err)
	}
	return nil
}

func DeleteGeneratedFiles(ctx context.Context, d *DirFile) error {
	files, err := os.ReadDir(d.Dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasPrefix(f.Name(), "output") || strings.HasSuffix(f.Name(), "_encoded.m3u8") || f.Name() == "list.txt" || f.Name() == "final_output.m3u8" {
			err = os.Remove(d.Dir + f.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
