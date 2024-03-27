package stablediffusion

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"tools/pkg/bos"
	"tools/pkg/config"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type StableDiffusion struct {
	ctx    context.Context
	Client http.Client
}

// Response 翻译结果
type Response struct {
	Images []string `json:"images"` // Images图像
}

var instance *StableDiffusion

var once sync.Once

var index int

func NewStableDiffusion(ctx context.Context) *StableDiffusion {
	once.Do(func() {
		instance = &StableDiffusion{
			ctx: ctx,
		}
	})
	return instance
}

// Image 获取生成结果
func (r *Response) Image() (image string, err error) {
	if len(r.Images) == 0 {
		err = errors.New("error text2img")
		return
	}
	image = r.Images[0]
	return
}

func (l *StableDiffusion) TextToImage(prompt string, negativePrompt string) (string, error) {

	strBody := l.buildBody(prompt, negativePrompt)

	hosts := config.GetString("sd.host")

	// 将hosts转换成数组，用|判断
	hostArr := strings.Split(hosts, "|")

	// 每次请求index 加1
	index = index + 1

	hostCount := len(hostArr)

	// 对index取余，得到host的索引
	hostIndex := index % hostCount
	url := "http://" + hostArr[hostIndex] + "/sdapi/v1/txt2img"
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(strBody))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := l.Client.Do(req)

	if err != nil {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "sd求错误",
			"body":  strBody,
			"error": err.Error(),
		})
		return "", err
	}

	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)

	// if err != nil {
	// 	wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
	// 		"type":     "error",
	// 		"msg":      "sd请求读取body错误",
	// 		"payload":  strBody,
	// 		"response": string(data),
	// 		"error":    err.Error(),
	// 	})
	// 	return "", err
	// }

	wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
		"type":    "info",
		"msg":     "sd请求成功",
		"payload": strBody,
		"data":    string(data),
	})

	res := &Response{}

	if err = json.Unmarshal(data, res); err != nil {

		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "大模型请响应解析失败",
			"data": string(data),
		})

		return "", err
	}

	imageStr, err := res.Image()

	if err != nil {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "大模型请响应失败",
			"data": string(data),
		})

	}

	imageStrDecode, err := base64.StdEncoding.DecodeString(imageStr)

	if err != nil {
		wailsruntime.EventsEmit(l.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "imageStrDecode",
			"data": string(data),
		})

	}

	imageStr = string(imageStrDecode)

	_, outURL, err := l.Upload(imageStr)

	return outURL, err
}

func (l *StableDiffusion) buildBody(prompt string, negativePrompt string) string {

	query := map[string]any{
		"width":           config.GetInt("sd.width"),
		"height":          config.GetInt("sd.height"),
		"prompt":          prompt,
		"negative_prompt": negativePrompt,
		"steps":           20,
		"cfg_scale":       7,
		"batch_size":      1,
		"seed":            -1,
		"sampler_index":   "DPM++ 2M Karras",
		"override_settings": map[string]string{
			"sd_model_checkpoint": "sdXL_v10",
			"sd_vae":              "sdxl_vae.safetensors",
		},
	}

	// 构建请求体
	body, _ := json.Marshal(query)

	return string(body)
}

// upload2Bos 上传到bos
func (l StableDiffusion) Upload(image string) (string, string, error) {

	client := bos.NewBos(l.ctx)

	return client.Upload("", image, "sd")
}
