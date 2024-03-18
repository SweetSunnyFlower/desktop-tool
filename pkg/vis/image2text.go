package vis

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"io"
	"net/http"
	"sync"
	"tools/pkg/config"

	uuid "github.com/satori/go.uuid"
)

type Vis struct {
	Client http.Client
	ctx    context.Context
}

type Payload struct {
	URI     string            `json:"uri"`
	Method  string            `json:"method"`
	Body    io.Reader         `json:"omitempty"`
	Headers map[string]string `json:"headers"`
	Cookies map[string]string `json:"cookies"`
}

type Response struct {
	SourceKey  string        `json:"source_key"`
	Status     string        `json:"status"`
	Code       int           `json:"code"`
	Logid      int64         `json:"logid"`
	FeatureRes FeatureResult `json:"feature_result"`
}

type FeatureResult struct {
	Feature    string `json:"feature"`
	Value      string `json:"value"`
	Status     string `json:"status"`
	CalcTimeMs int    `json:"calc_time_ms"`
}

type FeatureResultValue struct {
	ErrNo         int        `json:"err_no"`
	ErrMsg        string     `json:"err_msg"`
	Format        string     `json:"format"`
	Result        string     `json:"result"`
	FaceRet       string     `json:"face_ret"`
	HistoryMsg    [][]string `json:"history_msg"`
	OcrRet        string     `json:"ocr_ret"`
	WaitInQueueMs int        `json:"wait_in_queue_ms"`
	ClassifyMs    int        `json:"classify_ms"`
}

var VisInstalce *Vis

var once sync.Once

func NewVis(ctx context.Context) *Vis {
	once.Do(func() {
		wailsruntime.EventsEmit(ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "初始化Vis实例",
		})

		VisInstalce = &Vis{
			ctx: ctx,
		}
	})
	return VisInstalce
}

func (v *Vis) Image2Text(imageUrl string) (*FeatureResultValue, error) {
	imageData, err := v.GetImageData(imageUrl)

	if err != nil {
		return nil, err
	}

	resourceKey := v.Md5(base64.StdEncoding.EncodeToString(imageData))

	body, err := v.BuildBody(imageData, map[string]any{"session_id": resourceKey, "extra_info": ""})

	if err != nil {
		return nil, err
	}

	header := v.BuildHeader(imageData)

	url := config.GetString("vis.image2text_url") + "?business_name=" + config.GetString("vis.business_name") + "&feature_name=" + config.GetString("vis.feature_name")
	payload := &Payload{
		URI:     url,
		Method:  http.MethodPost,
		Headers: header,
		Body:    bytes.NewBuffer(body),
	}

	fullURL := fmt.Sprintf("%s%s", config.GetString("vis.host"), url)

	wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
		"type":    "info",
		"msg":     "图生文请求参数",
		"payload": payload,
		"fullURL": fullURL,
	})

	req, err := http.NewRequest(payload.Method, fullURL, payload.Body)

	// set 请求头
	for k, v := range payload.Headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	response, err := v.Client.Do(req)

	if err != nil {
		wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
			"type":    "error",
			"msg":     "图生文请求错误",
			"payload": payload,
			"fullURL": fullURL,
			"error":   err.Error(),
		})
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
			"type":  "error",
			"msg":   "图生文请求读取body错误",
			"data":  string(data),
			"error": err.Error(),
		})
		return nil, err
	}

	wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "图生文请求成功",
		"data": string(data),
	})

	var res interface{}

	// 这里返回的是框架返回的结果包裹业务返回的结果，生产环境都返回了，测试环境只返回了业务结果
	if config.GetBool("app.debug") {
		res = &FeatureResult{}
	} else {
		res = &Response{}
	}

	if err = json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	var featureResult *FeatureResult

	switch self := res.(type) {
	case *FeatureResult:
		if self.Value == "" {
			return nil, errors.New("vic feature value is empty")
		}
		featureResult = self
	case *Response:
		re, err := v.ParseFeatureResult(self)
		if err != nil {
			return nil, err
		}
		featureResult = re
	default:
		wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "图生文响应结果断言，未知类型",
		})
		return nil, errors.New("unknown type")
	}

	var featureResultValue *FeatureResultValue

	wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "图生文响应结果featureResult.Value",
		"data": featureResult.Value,
	})
	err = json.Unmarshal([]byte(featureResult.Value), &featureResultValue)

	if err != nil {
		wailsruntime.EventsEmit(v.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "图生文响应结果featureResult.Value解析失败",
			"data": err.Error(),
		})
		return nil, err
	}

	if featureResultValue == nil {
		return nil, errors.New("feature result value is empty")
	}

	return featureResultValue, nil
}

func (v *Vis) GetImageData(imageUrl string) ([]byte, error) {
	// 图片下载跳过https证书验证，2023.06.29 图片下载可用性下降
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (v *Vis) ParseFeatureResult(response *Response) (*FeatureResult, error) {

	if response.Code != 0 {
		return nil, errors.New("vic response code is not 0")
	}

	return &response.FeatureRes, nil

}

func (v *Vis) BuildBody(image []byte, otherParam map[string]any) ([]byte, error) {
	data := map[string]any{
		"image": base64.StdEncoding.EncodeToString(image),
	}
	for key, value := range otherParam {
		data[key] = value
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := map[string]any{
		"data": base64.StdEncoding.EncodeToString(dataJSON),
	}
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(body); err != nil {
		return nil, err
	}
	return bodyBytes, nil

}

func (v *Vis) Md5(data string) string {
	md5Byte := md5.Sum([]byte(data))
	return hex.EncodeToString(md5Byte[:])
}

func (v *Vis) BuildHeader(image []byte) map[string]string {
	resourceKey := v.Md5(base64.StdEncoding.EncodeToString(image))
	id := uuid.Must(uuid.NewV4(), nil)
	return map[string]string{
		"Content-Type":  "application/json; charset=UTF-8",
		"resource_key":  resourceKey,
		"auth_key":      config.GetString("vis.auth_key"),
		"business_name": config.GetString("vis.business_name"),
		"feature_name":  config.GetString("vis.feature_name"),
		"X_BD_LOGID":    id.String(),
	}
}
