package av1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	api_key = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxIiwianRpIjoiOTYxYWRmZmU5MzFkY2ViOGU0OTI3Y2M0Y2Q4YWIwZjI5ZWExOTgzNDJjZjc2MDYwNTU5OWY1OTUyYTE2ZmQwMWNhNGY3MzljNGIwZTZiOTQiLCJpYXQiOjE2ODYwNzAzNzAuMzM4Mzc5LCJuYmYiOjE2ODYwNzAzNzAuMzM4MzgsImV4cCI6NDg0MTc0Mzk3MC4zMzAxOSwic3ViIjoiNjM4ODI1MjQiLCJzY29wZXMiOlsicHJlc2V0LnJlYWQiLCJwcmVzZXQud3JpdGUiLCJ3ZWJob29rLndyaXRlIiwid2ViaG9vay5yZWFkIiwidGFzay53cml0ZSIsInRhc2sucmVhZCIsInVzZXIud3JpdGUiLCJ1c2VyLnJlYWQiXX0.FRWqdrIQZKeYQAOOfv0wXNg4RtWK79TdKONxnKNr34kCTlXwUJgqMiiy-BEAD-NAltJo1XCymbzO8kityeRUpTuHzcsZYmcgwod8NS1T8kOiyxIMRw7pBpiND-dULdFQFuVp2GeUGNGzdFYixas_Ft4KD87bN01BZBGbWKbUrbjQK7-PsxwSGuWAVwCW80_7WM1B-VL0Zo6JcAZ70cGeLUyTS_-9FdXKIsMFkuGNu1Y_8ldJH_J1oOVuaaSVN45ULOkSapznX2YPyFtF47Fp-JPmRdznqgYILMDpC8E5aGu_b7ked9m-sMVf4-Usf2AeVNB1uu2wiQzdR4BnRrW8k84EWMThNLuKie4O5qb4ovPctHtHJhcrN0OANGVtZ_bAJsCg1tRfaqWmu_5f3ncdlYYHXsC0qoBHs1yVxLAw0dgXc5K70We3g9EY4g20Fh8bC6tCjF-yb1df05zbdYbd1gIHTc-dwiOn6_WwMl5bF_4KnVJsbhqIAE2CCL0efoDoE2MMLJnh9jp7uJhf9v6GWOypo2w5WXhkyrulxGbbk9ma6IurSARONdICSyqs36ThdbrIoxO6lZMT05h8N8bK5wTji4_n19MEHxEZNbqE08Rv5UKkxzkqbtWjER6zGWMnU9F5hOi0nnxkH3aK6A5wY7IiJ2eFF3bqQnm4kH2JAnQ"
)

func ConvertToAV1(apiKey, inputPath, outputPath string) error {
	// 读取 M3U8 文件
	m3u8, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open M3U8 file: %v", err)
	}
	defer m3u8.Close()

	// 读取 M3U8 文件中的 TS 文件列表
	tsFiles, err := parseTSFiles(m3u8)
	if err != nil {
		return fmt.Errorf("failed to parse TS files: %v", err)
	}

	// 下载 TS 文件并保存到临时目录
	tempDir := inputPath[:strings.LastIndex(inputPath, "/")]

	// 将 TS 文件转换为 AV1 格式
	av1Filepath := outputPath
	av1File, err := os.Create(av1Filepath)
	if err != nil {
		return fmt.Errorf("failed to create AV1 file: %v", err)
	}
	defer av1File.Close()

	av1Options := map[string]interface{}{
		"video_codec":   "av1",
		"video_bitrate": 0,
		"audio_codec":   "none",
		"trim":          []interface{}{0, 0},
	}
	av1OptionsJSON, err := json.Marshal(av1Options)
	if err != nil {
		return fmt.Errorf("failed to encode AV1 options: %v", err)
	}

	requestBody := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(requestBody)

	// 添加 M3U8 文件
	m3u8Part, err := multipartWriter.CreateFormFile("input", filepath.Base(inputPath))
	if err != nil {
		return fmt.Errorf("failed to create M3U8 part: %v", err)
	}
	_, err = io.Copy(m3u8Part, m3u8)
	if err != nil {
		return fmt.Errorf("failed to write M3U8 part: %v", err)
	}

	// 添加 TS 文件
	for i, tsFile := range tsFiles {
		tsFilepath := filepath.Join(tempDir, tsFile)

		tsPart, err := multipartWriter.CreateFormFile(fmt.Sprintf("input_%d", i+1), tsFile)
		if err != nil {
			return fmt.Errorf("failed to create TS part #%d: %v", i, err)
		}

		tsFile, err := os.Open(tsFilepath)
		if err != nil {
			return fmt.Errorf("failed to open TS file #%d: %v", i, err)
		}
		defer tsFile.Close()

		_, err = io.Copy(tsPart, tsFile)
		if err != nil {
			return fmt.Errorf("failed to write TS part #%d: %v", i, err)
		}
	}

	// 添加转换选项
	optionsPart, err := multipartWriter.CreateFormField("options")
	if err != nil {
		return fmt.Errorf("failed to create options part: %v", err)
	}
	_, err = optionsPart.Write(av1OptionsJSON)
	if err != nil {
		return fmt.Errorf("failed to write options part: %v", err)
	}

	// 添加 API 密钥和转换器
	apiKeyPart, _ := multipartWriter.CreateFormField("apikey")
	apiKeyPart.Write([]byte(apiKey))

	converterPart, _ := multipartWriter.CreateFormField("converter")
	converterPart.Write([]byte("ffmpeg"))

	multipartWriter.Close()

	convertResp, err := http.Post("https://api.cloudconvert.com/v2/convert", multipartWriter.FormDataContentType(), requestBody)
	if err != nil {
		return fmt.Errorf("failed to request conversion: %v", err)
	}
	defer convertResp.Body.Close()

	convertRespBody, _ := ioutil.ReadAll(convertResp.Body)

	if convertResp.StatusCode != http.StatusOK {
		return fmt.Errorf("conversion failed with status code %d: %s", convertResp.StatusCode, convertRespBody)
	}

	// 下载转换后的 AV1 文件
	downloadURLResponse, _ := ioutil.ReadAll(convertResp.Body)

	downloadURLResponseJSON := struct {
		Data struct {
			Download struct {
				Url string `json:"url"`
			} `json:"download"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(downloadURLResponse, &downloadURLResponseJSON)
	if err != nil {
		return fmt.Errorf("failed to unmarshal download URL response: %v", err)
	}

	downloadURL := downloadURLResponseJSON.Data.Download.Url

	av1Resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download AV1 file: %v", err)
	}
	defer av1Resp.Body.Close()

	_, err = io.Copy(av1File, av1Resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save AV1 file: %v", err)
	}

	return nil
}

func parseTSFiles(m3u8 io.Reader) ([]string, error) {
	var tsFiles []string

	scanner := bufio.NewScanner(m3u8)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ".ts") {
			tsFiles = append(tsFiles, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tsFiles, nil
}

func ConvertMp4ToAv1(filePath string, outputPath string) {
	// 创建 REST 客户端
	client := resty.New()

	// 上传文件
	resp, err := client.R().
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", api_key)).
		SetBody(getFileContent(filePath)).
		Post("https://api.cloudconvert.com/v2/convert")

	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取任务 ID
	taskID, err := getTaskID(resp.Body())

	if err != nil {
		fmt.Println(err)
		return
	}

	// 等待转换完成
	for {
		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", api_key)).
			Get(fmt.Sprintf("https://api.cloudconvert.com/v2/tasks/%s", taskID))

		if err != nil {
			fmt.Println(err)
			return
		}

		status, err := getTaskStatus(resp.Body())

		if err != nil {
			fmt.Println(err)
			return
		}

		if status == "finished" {
			break
		}
	}

	// 下载转换后的文件
	resp, err = client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", api_key)).
		Get(fmt.Sprintf("https://api.cloudconvert.com/v2/tasks/%s/download", taskID))

	if err != nil {
		fmt.Println(err)
		return
	}

	// 将文件保存到本地
	err = ioutil.WriteFile(outputPath, resp.Body(), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Conversion completed successfully!")
}

// 解析任务 ID
func getTaskID(body []byte) (string, error) {
	// TODO: 解析 JSON，获取任务 ID
	return "", nil
}

// 解析任务状态
func getTaskStatus(body []byte) (string, error) {
	// TODO: 解析 JSON，获取任务状态
	return "", nil
}
