package av1

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ConvertTsToAv1(ctx context.Context, m3u8Path, tsDir string) (err error) {
	// 设置 CloudConvert API key 和上传文件路径
	//apiKey := api_key

	// 设置 OAuth2 客户端凭据
	//config := &clientcredentials.Config{
	//	ClientID:     "1145",
	//	ClientSecret: "MYfBuSg0TWZUret66bSUjLSzeILg8V2h7h74GVMF",
	//	TokenURL:     "https://api.cloudconvert.com/v2/auth/oauth/token",
	//}
	//
	//// 获取访问令牌
	//token, err := config.Token(context.Background())
	//if err != nil {
	//	fmt.Println("Failed to get access token:", err)
	//	return
	//}

	// 读取 M3U8 文件
	m3u8Data, err := ioutil.ReadFile(m3u8Path)
	if err != nil {
		fmt.Println("Failed to read M3U8 file:", err)
		return
	}

	// 读取 TS 文件并将它们添加到 M3U8 文件中
	tsFiles, err := filepath.Glob(tsDir + "/*.ts")
	if err != nil {
		fmt.Println("Failed to read TS files:", err)
		return
	}
	for _, tsFile := range tsFiles {
		m3u8Data = append(m3u8Data, []byte("\n"+tsFile)...)
	}

	// 将 M3U8 文件转换为 AV1 格式
	apiURL := "https://api.cloudconvert.com/v2/convert"
	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	// 添加 API key 和转换选项到请求体中
	//_ = writer.WriteField("apikey", apiKey)
	_ = writer.WriteField("input", "upload")
	_ = writer.WriteField("inputformat", "m3u8")
	_ = writer.WriteField("outputformat", "av1")
	_ = writer.WriteField("converteroptions", `{"video_codec": "av1"}`)

	// 添加 M3U8 文件到请求体中
	m3u8File, err := os.Open(m3u8Path)
	if err != nil {
		fmt.Println("Failed to open M3U8 file:", err)
		return
	}
	defer m3u8File.Close()
	m3u8Contents, err := ioutil.ReadAll(m3u8File)
	if err != nil {
		fmt.Println("Failed to read M3U8 file:", err)
		return
	}
	_ = writer.WriteField("file", base64.StdEncoding.EncodeToString(m3u8Contents))

	// 添加 TS 文件到请求体中
	var ts *os.File
	var tsContents []byte
	for _, tsFile := range tsFiles {
		ts, err = os.Open(tsFile)
		if err != nil {
			fmt.Println("Failed to open TS file:", err)
			return
		}
		defer ts.Close()
		tsContents, err = ioutil.ReadAll(ts)
		if err != nil {
			fmt.Println("Failed to read TS file:", err)
			return
		}
		var tsPart io.Writer
		tsPart, err = writer.CreateFormFile("files[]", filepath.Base(tsFile))
		if err != nil {
			fmt.Println("Failed to create form file:", err)
			return
		}
		_, _ = tsPart.Write(tsContents)
	}

	// 发送请求并获取响应
	_ = writer.Close()
	req, err := http.NewRequest(http.MethodPost, apiURL, requestBody)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+api_key)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}
	fmt.Println(string(response))
	return nil
}

// 获取文件内容
func getFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

// 获取下载链接
//func getDownloadURL(client *resty.Client, taskID string) string {
//	resp, err := client.R().
//		SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey)).
//		Get(fmt.Sprintf("https://api.cloudconvert.com/v2/tasks/%s/download", taskID))
//
//	if err != nil {
//		panic(err)
//	}
//
//	// 解析响应，获取下载链接
//	downloadURL := resp.Body().GetObject()["data"].GetObject()["url"].GetString()
//	return downloadURL
//}

// 下载文件
func downloadFile(url string, filePath string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}
