package vis

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"net/http"
	"sync"
	"tools/pkg/config"

	uuid "github.com/satori/go.uuid"
)

type Vis struct {
	Client http.Client
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
	ErrNo         int    `json:"err_no"`
	ErrMsg        string `json:"err_msg"`
	Format        string `json:"format"`
	Result        string `json:"result"`
	WaitInQueueMs int    `json:"wait_in_queue_ms"`
	ClassifyMs    int    `json:"classify_ms"`
}

var VisInstalce *Vis

var once sync.Once

func NewVis() *Vis {
	once.Do(func() {
		VisInstalce = &Vis{}
	})
	return VisInstalce
}

func (v *Vis) Image2Text(imageUrl string) (string, error) {
	imageData := []byte("")

	resourceKey := v.Md5(base64.StdEncoding.EncodeToString(imageData))

	body, err := v.BuildBody(imageData, map[string]any{"session_id": resourceKey, "extra_info": ""})

	if err != nil {
		return "", err
	}

	header := v.BuildHeader(imageData)

	url := config.GetString("vis.image2text_url") + "?business_name=" + config.GetString("app.business_name") + "&feature_name=" + config.GetString("app.feature_name")
	payload := &Payload{
		URI:     url,
		Method:  http.MethodPost,
		Headers: header,
		Body:    bytes.NewBuffer(body),
	}

	fullURL := fmt.Sprintf("%s/%s", config.GetString("vis.host"), url)

	req, err := http.NewRequest(payload.Method, fullURL, payload.Body)

	if err != nil {
		return "", err
	}

	response, err := v.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	var res any

	// 这里返回的是框架返回的结果包裹业务返回的结果，生产环境都返回了，测试环境只返回了业务结果
	if config.GetBool("app.debug") {
		res = &FeatureResult{}
	} else {
		res = &Response{}
	}

	if err = json.Unmarshal(data, res); err != nil {
		return "", err
	}

	// 生产环境，需要解读出来业务结果
	if !config.GetBool("app.debug") {
		res, err = v.ParseFeatureResult(res.(*Response))
		if err != nil {
			return "", err
		}
	}

	if res.(*FeatureResult).Value == "" {
		return "", errors.New("vic feature value is empty")
	}

	featureResultValue := FeatureResultValue{}
	_ = json.Unmarshal([]byte(res.(*FeatureResult).Value), &featureResultValue)

	if featureResultValue.Result == "" {
		return "", errors.New("feature result value is empty")
	}

	return featureResultValue.Result, nil
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
