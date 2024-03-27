package translate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"tools/pkg/config"

	"github.com/pkg/errors"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Result 翻译结果
type Result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// Response 翻译结果
type Response struct {
	Errno   int      `json:"error_code"`
	Error   string   `json:"error_msg"`
	Results []Result `json:"trans_result"`
}

type BaiduTranslate struct {
	Client http.Client
	ctx    context.Context
}

// Dst 获取翻译结果
func (r *Response) Dst() (dst string, err error) {
	if r.Errno != 0 {
		err = errors.Errorf("error: error_code = %d error_msg = %s", r.Errno, r.Error)
		return
	}
	if len(r.Results) == 0 {
		err = errors.Errorf("error trans result empty: error_code = %d error_msg = %s", r.Errno, r.Error)
		return
	}
	dst = r.Results[0].Dst
	return
}

var BaiduTranslateInstance *BaiduTranslate

var once sync.Once

func NewBaiduTranslate(ctx context.Context) *BaiduTranslate {
	once.Do(func() {
		wailsruntime.EventsEmit(ctx, "logEvent", map[string]interface{}{
			"type": "info",
			"msg":  "初始化翻译实例",
		})

		BaiduTranslateInstance = &BaiduTranslate{
			ctx: ctx,
		}
	})
	return BaiduTranslateInstance
}

func (t *BaiduTranslate) Translate(text string, fromLang, toLang string) (string, error) {
	body := t.BuildBody(text, fromLang, toLang)

	fullURL := fmt.Sprintf("%s%s", config.GetString("translate.host"), "/api/trans/vip/translate")

	wailsruntime.EventsEmit(t.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "翻译body:" + body,
	})

	req, err := http.NewRequest(http.MethodPost, fullURL, strings.NewReader(body))

	if err != nil {
		return "", err
	}

	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	response, err := t.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer func() {
		_ = response.Body.Close()
	}()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		wailsruntime.EventsEmit(t.ctx, "logEvent", map[string]interface{}{
			"type": "error",
			"msg":  "翻译失败:" + err.Error(),
		})
		return "", err
	}
	wailsruntime.EventsEmit(t.ctx, "logEvent", map[string]interface{}{
		"type": "info",
		"msg":  "响应成功:" + string(data),
	})
	translateResponse := &Response{}
	if err = json.Unmarshal(data, translateResponse); err != nil {
		return "", err
	}

	return translateResponse.Dst()
}

func (t *BaiduTranslate) BuildBody(text string, from string, to string) string {
	appID := config.GetString("translate.app_id")
	salt := "0000000"
	md516 := md5.Sum([]byte(appID + text + salt + config.GetString("translate.app_secret")))
	sign := hex.EncodeToString(md516[:])
	var v = url.Values{}
	v.Add("q", text)
	v.Add("from", from)
	v.Add("to", to)
	v.Add("appid", appID)
	v.Add("salt", salt)
	v.Add("sign", sign)
	return v.Encode()
}
